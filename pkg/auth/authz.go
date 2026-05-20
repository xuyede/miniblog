// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// casbin访问控制器
	aclModel = `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
	`
)

type Authz struct {
	*casbin.SyncedEnforcer
}

// NewAuthz 创建一个使用 casbin 完成授权的授权器.
func NewAuthz(db *gorm.DB) (*Authz, error) {
	// 用你已有的 *gorm.DB 创建 Casbin 的策略存储适配器
	// 它会自动在数据库中创建 casbin_rule 表（如果不存在）
	// 策略数据（谁能访问什么资源）存在这张表里
	adapter, err := adapter.NewAdapterByDB(db)

	if err != nil {
		return nil, err
	}

	// 加载访问控制模型
	m, _ := model.NewModelFromString(aclModel)

	// 创建一个同步的访问控制器，使用上面创建的适配器来存储策略数据
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 从数据库加载权限规则到内存中
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	// 每 5 秒从数据库重新加载策略，这样你在数据库里改了权限规则后，不需要重启服务就能生效
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{enforcer}

	return a, nil
}

// Authorize 用来进行授权.
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}
