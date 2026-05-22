// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package core

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/xuyede/miniblog/internal/pkg/errno"
)

// BindAndValidate 从请求 body 中反序列化 JSON 并使用 govalidator 校验字段合法性。
// 成功返回解析后的结构体指针，失败则自动写入错误响应并返回 nil。
func BindAndValidate[T any](c *gin.Context) *T {
	var r T

	// 如果是GET请求，应该直接从请求链接中获取参数，否则从请求 body 中获取参数
	if c.Request.Method == "GET" {
		if err := c.ShouldBindQuery(&r); err != nil {
			GenarateResponse(c, errno.ErrBind, nil)
			return nil
		}
	} else {
		if err := c.ShouldBindJSON(&r); err != nil {
			GenarateResponse(c, errno.ErrBind, nil)
			return nil
		}
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		GenarateResponse(c, errno.ErrInvalidParameter.SetMessage("%s", err.Error()), nil)
		return nil
	}

	return &r
}
