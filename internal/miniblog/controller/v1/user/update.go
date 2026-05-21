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

func (ctrl *UserController) Update(c *gin.Context) {
	log.C(c).Infow("PUT /v1/users/:name called")

	r := core.BindAndValidate[v1.UpdateUserRequest](c)
	if r == nil {
		return
	}

	if err := ctrl.b.Users().Update(c, c.Param("name"), r); err != nil {
		core.GenarateResponse(c, err, nil)
		return
	}

	core.GenarateResponse(c, nil, nil)
}
