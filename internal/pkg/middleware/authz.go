// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xuyede/miniblog/internal/pkg/core"
	"github.com/xuyede/miniblog/internal/pkg/errno"
	"github.com/xuyede/miniblog/internal/pkg/known"
	"github.com/xuyede/miniblog/internal/pkg/log"
)

type Auther interface {
	Authorize(sub, obj, act string) (bool, error)
}

func Authz(a Auther) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString(known.XUsernameKey)
		obj := c.Request.URL.Path
		act := c.Request.Method

		log.Debugw("授权上下文", "sub", sub, "obj", obj, "act", act)

		// 接口解耦：中间件依赖 Auther 接口（只要有 Authorize 方法），不直接依赖 Casbin 具体实现
		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			core.GenarateResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
	}
}
