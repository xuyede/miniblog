// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"

	"github.com/xuyede/miniblog/internal/miniblog/store"
	"github.com/xuyede/miniblog/internal/pkg/errno"
	"github.com/xuyede/miniblog/internal/pkg/model"
	v1 "github.com/xuyede/miniblog/pkg/api/miniblog/v1"
)

// UserBiz 接口的实现.
type userBiz struct {
	ds store.IStore
}

// New 创建一个实现了 UserBiz 接口的实例.
func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

// UserBiz 定义了 user 模块在 biz 层所实现的方法.
type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

func (u *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := u.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)
