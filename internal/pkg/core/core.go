// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuyede/miniblog/internal/pkg/errno"
)

type ErrResponse struct {
	// Code 指定了业务错误码
	Code string `json:"code"`

	// Message 包含了可以直接对外展示的错误信息
	Message string `json:"msg"`
}

type SuccessResponse struct {
	// Code 指定了业务错误码，成功时通常为 200
	Code int `json:"code"`

	// Message 包含了可以直接对外展示的信息，成功时通常为 "success"
	Message string `json:"msg"`

	// Data 包含了响应数据，成功时可以包含实际的响应内容
	Data interface{} `json:"data"`
}

// 将错误或响应数据写入 HTTP 响应主体。
// 使用 errno.Decode 方法，根据错误类型，尝试从 err 中提取业务错误码和错误信息.
func GenarateResponse(c *gin.Context, err error, data interface{}) {

	if err != nil {
		httpCode, code, message := errno.Decode(err)
		c.JSON(httpCode, ErrResponse{
			Code:    code,
			Message: message,
		})
		return
	}

	// 如果 err 为 nil，说明没有错误发生，直接返回成功响应和数据
	c.JSON(http.StatusOK, SuccessResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}
