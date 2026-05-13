// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package verflag

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/xuyede/miniblog/pkg/version"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const (
	strRawVersion   = "raw"
	versionFlagName = "version"
)

// versionFlag 是一个全局变量，表示 `--version` 标志的值.
var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

// Version 包装了 VersionVar 函数.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)

	return p
}

// VersionVar 定义了一个具有指定名称和用法的标志.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	// 把设置 --version 的操作指针给了 p （就是 *versionFlag）
	pflag.Var(p, name, usage)
	// `--version` 等价于 `--version=true`
	// 用户传了 --version（无值） → NoOptDefVal 生效 → pflag 调用 versionFlag.Set("true") → *versionFlag = VersionTrue（1）
	pflag.Lookup(name).NoOptDefVal = "true"
}

// AddFlags 将 `--version` 标志添加到指定的 FlagSet 中.
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// PrintAndExitIfRequested 将检查是否传递了 `--version` 标志，如果是，则打印版本并退出.
func PrintAndExitIfRequested() {
	switch *versionFlag {
	case VersionRaw:
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	case VersionTrue:
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

// String 实现了 pflag.Value 接口中的 String 方法.
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// Set 实现了 pflag.Value 接口中的 Set 方法.
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw

		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

// Type 实现了 pflag.Value 接口中的 Type 方法.
func (v *versionValue) Type() string {
	return "version"
}
