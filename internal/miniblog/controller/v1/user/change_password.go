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

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("POST /v1/users/:name/change_password called")

	r := core.BindAndValidate[v1.ChangePasswordRequest](c)
	if r == nil {
		return
	}

	if err := ctrl.b.Users().ChangePassword(c, c.Param("name"), r); err != nil {
		core.GenarateResponse(c, err, nil)

		return
	}

	core.GenarateResponse(c, nil, nil)
}
