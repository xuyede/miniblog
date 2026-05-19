// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package user

import (
	"github.com/gin-gonic/gin"

	"github.com/xuyede/miniblog/internal/pkg/core"
	"github.com/xuyede/miniblog/internal/pkg/errno"
	"github.com/xuyede/miniblog/internal/pkg/log"
	v1 "github.com/xuyede/miniblog/pkg/api/miniblog/v1"
)

func (ctrl *UserController) Login(c *gin.Context) {
	log.C(c).Infow("POST /login called")

	var r v1.LoginRequest
	// 把请求 body 反序列化为 LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.GenarateResponse(c, errno.ErrBind, nil)
		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.GenarateResponse(c, err, nil)
		return
	}

	core.GenarateResponse(c, nil, resp)
}
