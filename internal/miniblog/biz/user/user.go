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
	"github.com/xuyede/miniblog/pkg/auth"
	"github.com/xuyede/miniblog/pkg/token"
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
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error
}

func (u *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM

	// 把 CreateUserRequest 的字段拷贝到 model.UserM 中.
	_ = copier.Copy(&userM, r)

	if err := u.ds.Users().Create(ctx, &userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}

		return err
	}

	return nil
}

func (u *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 获取登录用户的所有信息
	user, err := u.ds.Users().Get(ctx, r.Username)

	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

func (u *userBiz) ChangePassword(ctx context.Context, username string, r *v1.ChangePasswordRequest) error {
	// 获取用户的所有信息
	user, err := u.ds.Users().Get(ctx, username)
	if err != nil {
		return errno.ErrUserNotFound
	}

	// 对比传入的旧密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	// 更新用户密码
	user.Password, _ = auth.Encrypt(r.NewPassword)
	if err := u.ds.Users().Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)
