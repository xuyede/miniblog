// Copyright 2026 Deye <xuyede77@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/xuyede/miniblog.

package miniblog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/xuyede/miniblog/internal/miniblog/controller/v1/user"
	"github.com/xuyede/miniblog/internal/miniblog/store"
	"github.com/xuyede/miniblog/internal/pkg/known"
	"github.com/xuyede/miniblog/internal/pkg/log"
	mw "github.com/xuyede/miniblog/internal/pkg/middleware"
	pb "github.com/xuyede/miniblog/pkg/proto/miniblog/v1"
	"github.com/xuyede/miniblog/pkg/token"
	"github.com/xuyede/miniblog/pkg/version/verflag"
)

var cfgFile string

// NewMiniBlogCommand 创建一个 *cobra.Command 对象. 之后，可以使用 Command 对象的 Execute 方法来启动应用程序.
func NewMiniBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 指定命令的名字，该名字会出现在帮助信息中
		Use: "miniblog",
		// 命令的简短描述
		Short: "A good Go practical project",
		// 命令的详细描述
		Long: `A good Go practical project, used to create user with basic information.

Find more miniblog information at:
        https://github.com/xuyede/miniblog#readme`,

		// 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果 `--version=true`，则打印版本并退出
			verflag.PrintAndExitIfRequested()

			// 初始化日志
			log.Init(logOptions())
			defer log.Sync()
			return run()
		},
		// 这里设置命令运行时，不需要指定命令行参数
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)

	// Cobra 支持持久性flag(PersistentFlag)，该flag可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the miniblog configuration file. Empty string for no configuration file.")

	// Cobra 也支持本地flag，本地flag只能在其所绑定的命令上使用
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 注册 version flag
	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

// run 函数是实际的业务代码入口函数.
func run() error {
	// 初始化 store 层
	if err := initStore(); err != nil {
		log.Errorw("初始化 store 层失败", "err", err)
		return err
	}

	// 设置 token 包的签发密钥，用于 token 包 token 的签发和解析
	token.Init(viper.GetString("jwt-secret"), known.XUsernameKey)

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("runmode"))

	// 创建一个不带任何中间件的路由引擎
	g := gin.New()

	// 注册中间件
	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestID()}
	g.Use(mws...)

	if err := installRouters(g); err != nil {
		return err
	}

	// 创建并运行 HTTP 服务器
	httpsrv := startInsecureServer(g)

	// 创建并运行 HTTPS 服务器
	httpssrv := startSecureServer(g)

	// 创建并运行 GRPC 服务器
	grpcsrv := startGRPCServer()

	// 等待中断信号优雅地关闭服务器（10 秒超时)
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 CTRL + C 就是触发系统 SIGINT 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 管道取值阻塞，等待关闭信号后再执行
	<-quit
	log.Infow("正在关闭服务...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP 服务强制退出", "err", err)
		return err
	}
	if err := httpssrv.Shutdown(ctx); err != nil {
		log.Errorw("HTTPS 服务强制退出", "err", err)
		return err
	}

	grpcsrv.GracefulStop()

	log.Infow("服务已退出")

	return nil
}

func startSecureServer(g *gin.Engine) *http.Server {
	httpsSrv := &http.Server{
		Addr:    viper.GetString("tls.addr"),
		Handler: g,
	}

	log.Infow("监听 TLS 服务端口", "addr", viper.GetString("tls.addr"))
	cert, key := viper.GetString("tls.cert"), viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			if err := httpsSrv.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalw("启动 HTTPS Server 失败", "error", err.Error())
			}
		}()
	}

	return httpsSrv
}

func startInsecureServer(g *gin.Engine) *http.Server {
	// 创建 HTTP Server 实例
	httpSrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	log.Infow("监听服务端口", "addr", viper.GetString("addr"))
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw("启动 HTTP Server 失败", "error", err.Error())
		}
	}()

	return httpSrv
}

func startGRPCServer() *grpc.Server {
	// 监听 GRPC 服务端口
	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
	}

	// 创建 GRPC Server 实例
	grpcsrv := grpc.NewServer()
	// user.New(store.S, nil)执行返回 *UserController 对象
	// 使得 GRPC Server 能够通过 *UserController 的接口去处理来自客户端的请求
	pb.RegisterMiniBlogServer(grpcsrv, user.New(store.S, nil))

	// 运行 GRPC 服务器。在 goroutine 中启动服务器，它不会阻止下面的正常关闭处理流程
	// 打印一条日志，用来提示 GRPC 服务已经起来，方便排障
	log.Infow("监听 GRPC 服务端口", "addr", viper.GetString("grpc.addr"))
	go func() {
		if err := grpcsrv.Serve(lis); err != nil {
			log.Fatalw("启动 GRPC Server 失败", "error", err.Error())
		}
	}()

	return grpcsrv
}
