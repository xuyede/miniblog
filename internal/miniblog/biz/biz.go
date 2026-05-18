// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package biz

import (
	"github.com/xuyede/miniblog/internal/miniblog/biz/user"
	"github.com/xuyede/miniblog/internal/miniblog/store"
)

type biz struct {
	ds store.IStore
}

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

type IBiz interface {
	Users() user.UserBiz
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}

// 确保 biz 实现了 IBiz 接口.
// 也可以写成 var _ IBiz = &biz{}
var _ IBiz = (*biz)(nil)
