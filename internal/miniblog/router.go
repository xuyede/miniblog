// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package miniblog

import (
	"github.com/gin-gonic/gin"
	"github.com/xuyede/miniblog/internal/miniblog/controller/v1/user"
	"github.com/xuyede/miniblog/internal/miniblog/store"
	"github.com/xuyede/miniblog/internal/pkg/core"
	"github.com/xuyede/miniblog/internal/pkg/errno"
	"github.com/xuyede/miniblog/internal/pkg/log"
	mw "github.com/xuyede/miniblog/internal/pkg/middleware"
)

func installRouters(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.GenarateResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.GenarateResponse(c, nil, map[string]string{"status": "ok"})
	})

	// 创建 user 模块的 controller，并将 store.S 传递给它，以创建 Store .
	uc := user.New(store.S)

	// 注册 /login
	g.POST("/login", uc.Login)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)

			// 该路由组需要认证中间件进行保护，只有合法用户才能访问
			userv1.Use(mw.Authn())
			userv1.PUT(":name/change-password", uc.ChangePassword)
		}
	}

	return nil
}
