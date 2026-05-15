// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package errno

import "fmt"

// Errno 包含 HTTP 状态码、业务错误码和错误信息
type Errno struct {
	HTTP    int
	Code    string
	Message string
}

// 实现 error 接口的 Error 方法，使 Errno 类型满足 error 接口的要求. 这样我们就可以将 Errno 类型的值作为 error 类型来使用了.
func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(format string, args ...interface{}) *Errno {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// Decode 将 error 类型转换为 Errno 类型，如果无法转换，则返回 InternalServerError.
func Decode(err error) (int, string, string) {
	switch types := err.(type) {
	case *Errno:
		return types.HTTP, types.Code, types.Message
	default:
	}

	// 默认返回未知错误码和错误信息. 该错误代表服务端出错
	return InternalServerError.HTTP, InternalServerError.Code, InternalServerError.Message
}
