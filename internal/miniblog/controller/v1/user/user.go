// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package user

import (
	"github.com/xuyede/miniblog/internal/miniblog/biz"
	"github.com/xuyede/miniblog/internal/miniblog/store"
	"github.com/xuyede/miniblog/pkg/auth"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求.
type UserController struct {
	b biz.IBiz
	a *auth.Authz
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{
		a: a,
		b: biz.NewBiz(ds),
	}
}
