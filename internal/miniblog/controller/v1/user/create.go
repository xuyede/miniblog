// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package user

import (
	"github.com/gin-gonic/gin"

	"github.com/xuyede/miniblog/internal/pkg/core"
	"github.com/xuyede/miniblog/internal/pkg/log"
	v1 "github.com/xuyede/miniblog/pkg/api/miniblog/v1"
)

const defaultMethods = "(GET)|(POST)|(PUT)|(DELETE)"

func (ctrl *UserController) Create(c *gin.Context) {
	log.C(c).Infow("POST /v1/users called")

	r := core.BindAndValidate[v1.CreateUserRequest](c)
	if r == nil {
		return
	}

	if err := ctrl.b.Users().Create(c, r); err != nil {
		core.GenarateResponse(c, err, nil)

		return
	}

	// 创建用户成功后，给该用户添加访问 /v1/users/:name 的权限
	if _, err := ctrl.a.AddNamedPolicy("p", r.Username, "/v1/users/"+r.Username, defaultMethods); err != nil {
		core.GenarateResponse(c, err, nil)
		return
	}

	core.GenarateResponse(c, nil, nil)
}
