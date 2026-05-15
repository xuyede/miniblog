## 一、初始化项目

### 1. 初始化GO模块
- go mod init <project>

### 2. 初始化Git
- 使用 [.gitignore](https://www.toptal.com/developers/gitignore)快速生成相关项目的模版

### 3. 热加载工具

- 使用 [air](https://github.com/air-verse/air/blob/master/README-zh_cn.md)

```go install github.com/cosmtrek/air@latest```

**配置 .air.toml 文件**

```执行：air -c .air.toml```

### 4. 使用Swagger编写接口文档
1. ```go install github.com/go-swagger/go-swagger/cmd/swagger@latest```

2. ```swagger serve -F=swagger --no-open --port 65534 ./api/openapi/openapi.yaml```

3. 在浏览器中打开 http://localhost:65534/docs

### 5. 添加 LICENSE
1. ```go install github.com/nishanths/license/v5@latest```

2. ```license -n 'Deye(许业德) <xuyede77@gmail.com>' -o LICENSE mit```

### 6. 给源文件添加版本声明

1. ```go install github.com/marmotedu/addlicense@latest```

2. ```addlicense -v -f ./scripts/boilerplate.txt --skip-dirs=third_party,vendor,_output .cmd/miniblog/main.go```

### 7. 使用Makefile实现上面的步骤
**配置 Makefile 文件**

- Makefile [学习](https://github.com/marmotedu/geekbang-go/blob/master/makefile/Makefile%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86.md)

## 二、应用组成及构建

### 1. 引用配置

- 命令行选项、命令行参数： 选择 [pflag](https://github.com/spf13/pflag)；

- 配置文件： 建议选择支持多种配置文件格式的包，也即从：viper、configor、koanf、config 中选择其一，毫无疑问 [viper](https://github.com/spf13/viper) 胜出：
  - viper 有 21000+ 的 Star 数，比其它包更受欢迎；
  - viper 功能强大，并且经过很多大型项目验证过；
  - viper 同时也可以实现分布式配置中心的功能。

- 环境变量： 如果环境变量不多，可以使用 os.Getenv，如果环境变量很多，可以使用 [envconfig](https://github.com/kelseyhightower/envconfig) 直接将环境变量读取到 Go 结构体变量中。

- 配置中心： 可根据需要选择 viper、apollo、[etcd](https://pkg.go.dev/go.etcd.io/etcd/client/v3)、[consul](https://github.com/hashicorp/consul) 等。其实一般的项目不需要引入配置中心，因为使用配置中心，会带来一些部署、维护的复杂度。

**结论：直接使用`pflag` + `viper`**

### 2. 应用业务逻辑

应用的业务逻辑根据业务的不同差别很大。一般而言，一个 Go 应用中会执行以下类别的业务逻辑处理（可能会用到其中一个或多个）：

- 初始化缓存；

- 初始化并创建各类数据库客户端，例如：Redis、MySQL、Kafka、MongoDB、Etcd 等；

- 初始化并创建其他服务的客户端等；

- 初始化并启动Web服务，例如：HTTP、HTTPS、GRPC；

- 启动异步任务，这些异步任务可以执行任何业务需要的操作，例如：watch kube-apiserver、定期从第三方服务拉取数据，并缓存、注册 /metrics 并监听指定的端口、启动 kafka 消费队列等等；

- 执行特定的业务处理，并退出程序；

### 3. 应用启动框架

启动框架你可以理解为一个 main 函数，只不过这里的 main 函数是有代码结构的，并可能分散在多个 Go 源码文件中，在这个大函数中，你可以读取配置文件、初始化业务逻辑、启动 Web 服务等，例如

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/spf13/pflag"
)

const helpText = `Usage: main [flags] arg [arg...]

This is a very simple app framework (does nothing).

Flags:`

var (
    addr = pflag.String("addr", ":8777", "The address to listen to.")
    help = pflag.BoolP("help", "h", false, "Show this help message.")

    usage = func() {
        fmt.Println(helpText)
        pflag.PrintDefaults()
    }
)

func main() {
    // 1. 命令行参数处理：解析，并读取命令行参数
    pflag.Usage = usage
    pflag.Parse()
    if *help {
        pflag.Usage()
        return
    }

    // 2. 业务处理：初始化路由
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello world")
    })
    server := http.Server{Addr: *addr}
    fmt.Printf("Starting http server at %s\n", *addr)

    // 3. 业务处理：启动 HTTP Web 服务
    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}
```

业务简单可以采用平铺式写法，但是多了就需要工具辅助，推荐[cobra](https://github.com/spf13/cobra)

### 4. 最佳构建方法

**结论: 使用 `pflag`、`viper`、`cobra` 来构建一个强大的应用程序**

- pflag - 给应用添加命令行标识 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8Pflag%E7%BB%99%E5%BA%94%E7%94%A8%E6%B7%BB%E5%8A%A0%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A0%87%E8%AF%86.md)
- viper - 获取配置文件数据 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90%E7%A5%9E%E5%99%A8-Viper%E5%85%A8%E8%A7%A3.md)
- cobra - 命令行框架 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E7%8E%B0%E4%BB%A3%E5%8C%96%E7%9A%84%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A1%86%E6%9E%B6-Cobra%E5%85%A8%E8%A7%A3.md)

### 5. 选择配置读取

**使用 `.yaml` 格式的配置文件来配置应用，并使用 `viper` 读取配置**

- 查看 `miniblog\internal\miniblog\helper.go`


## 三、设计日志包

日志包有很多，优先考虑 [zap](https://liwenzhou.com/posts/go/zap/) 和 [logrus](https://www.cnblogs.com/binHome/p/12027471.html)

- 如果对性能要求不高，追求使用简单，可以选择 logrus

- 如果对性能要求较高，并且追求相对方便的使用方式，可以选择 zap

项目使用 `zap` 作为日志包

- 快速学习 - [zap](https://github.com/marmotedu/geekbang-go/blob/master/%E4%BC%98%E7%A7%80%E5%BC%80%E6%BA%90%E6%97%A5%E5%BF%97%E5%8C%85%E4%BD%BF%E7%94%A8%E6%95%99%E7%A8%8B.md#zap%E5%8C%85%E4%BB%8B%E7%BB%8D)

### 1. 实现miniblog的logger

**参考 `miniblog\internal\pkg\log.go`**

**从已有的代码入手，优先在github社区找，查找方法：**

查找1：尝试在 zap 官方仓库的 README 文件中和仓库中的 examples 这类目录中查找看是否有创建示例（官方仓库是最可能存放这种示例代码的地方。因为是官方仓库，所代码质量会比较高，所以要优先从官方仓库中找）

查找2：尝试在 GitHub 上找
![alt text](image-1.png)
- 在 GitHub 搜索栏，输入 language:go zap demo 以搜索仓库名/仓库描述中同时有 zap 和 demo 关键字的 Go 代码仓库
- 按 Most stars 排序
- 从上到下，阅读检索出的代码仓库，根据代码仓库名和描述，判断是否是可以参考的 Go 项目，如果是，则进入仓库进行更详细的了解

查找3：进行更深入的查找，查找使用了 zap 包的代码，根据代码来判断 Go 代码段或者代码段所在的项目是否可以借鉴使用
![alt text](image-2.png)
- 在 GitHub 搜索栏，输入 language:go go.uber.org/zap/zapcore 以搜索可能封装了 zap 包的代码段。查看搜索结果的Code
- 根据代码段内容判断该代码段是否是可能的参考对象。如果是，则打开文件阅读源码，如果觉得源码可以利用，可以再进一步了解其所在的代码仓库，也许你会发现这个代码仓库就是一个完整的实现。

## 四、给应用添加版本信息，方便排查

添加一个`--version`的命令，触发时输出版本信息，git信息

**查看`miniblog\pkg\version`**

## 五、开发Web服务

Go Web 服务中最常用的 API 风格是：`REST` 和 `RPC`。`REST 底层使用的是 HTTP 协议`，`RPC 底层使用的是 RPC 协议`。每种协议，又有其适配的数据交换格式。这种适配，你可以理解为一种事实上的标准，如无特殊需求，无需打破这种适配关系：`REST API 风格采用 JSON 数据格式`，`RPC API 风格采用 Protobuf 数据格式`。

REST 和 RPC 又有其适配的场景，在企业应用开发中，通常会采用 2 种通信协议，一起构建一个优秀的 Go 应用。

- 对外（对接客户端）： REST + JSON 的组合。因为 API 接口规范、数据格式直观、易懂、开发调试方法，再加上客户端和服务端通过 HTTP 协议通信时，无需使用相同的编程语言，所以 REST + JSON 更适合对外提供 API 接口；

- 对内（服务端内部请求）： RPC + Protobuf 的组合。因为 RPC 协议调用方便，Protobuf 格式数据传输效率更高，从而使得 API 接口性能更好，所以 RPC + Probobuf 的组合更适合对内提供接口。

### 1. http请求处理流程
![alt text](image.png)

### 2. 实现一个简单的 REST Web Server

**选择[Gin](https://github.com/gin-gonic/gin)作为 REST Web框架，可以通过[示范](https://github.com/gin-gonic/examples)学习Gin的使用。一句话，Gin 帮你屏蔽了底层的路由匹配、请求解析、响应序列化等重复工作，你只需要关注 handler 里业务逻辑**

一个简单的例子：
```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

### 3. 使用Gin开发服务

**搜索关键字`gin.New()`**

### 4. 服务中间件 Middleware

在 Web 开发中，我们要实现很多功能，例如：`认证、授权、限流、熔断、设置请求/返回 Header（例如：请求 ID）、跨域`等，这些都需要通过 Web 中间件的方式来实现，可以说中间件是 Web 框架或者 Web 服务非常核心的一个功能，一个中大型的 Web 应用基本都需要用到。

那么什么是 Web 中间件呢？简单来说，Web 中间件是 HTTP / RPC 请求必经的一个中间层，该中间层可以统一处理所有的请求，你可以根据需要开发不同功能的中间层。例如：你可以在中间层给所有的请求/返回头中添加 `X-Request-ID` 头，用以标识唯一一次请求，方便追踪、排障。

基本上所有的 Web 框架都具有中间件的能力，不同框架可能叫法不同，例如有叫 Filter、Middleware 的，但实现的都是类似的机制。中间件常用在权限验证、日志记录、数据过滤等场景中。

### 5. 中间件在 Gin 中的实现

**Gin相关的中间件库：[点这](https://github.com/gin-gonic/contrib)**

Gin 也具有强大的中间件能力，Gin 的中间件是基于洋葱模型的，如下图所示：
![alt text](image-3.png)

上图中，有 2 个中间件：Middleware A、Middleware B。HTTP 请求，从开始到结束经历的路径为：Middleware A -> Middleware B -> 主体函数 -> Middleware B -> Middleware A。`执行顺序类似于栈`。

从上图可以知道，Gin 中间件，其实可以起到请求前置拦截和后置拦截的功能：

- 请求前置拦截： Web 请求到达我们定义的 HTTP 请求处理方法之前，拦截请求并进行相应处理；

- 请求后置拦截： 在处理完成请求并响应客户端时，拦截响应并进行相应的处理。

### 6. 中间件类型

#### 全局中间件：全局中间件设置之后对全局的路由都起作用

`r.Use()`，可以根据需要设置 1 个或者多个：

```go
r := gin.New()
//一次设置多个中间件
r.Use(Logger(), Recovery())
//一次设置一个中间件
r.Use(gin.Logger())
r.Use(gin.Recovery())
```

#### 路由组中间件：路由组中间件仅对该路由组下面的路由起作用

`r.Droup()`

```go
r := gin.New()

// 声明同时设置
authorized1 := r.Group("/users", AuthRequired())

// 先声明路由，再通过User添加
authorized2 := r.Group("/users")
authorized2.Use(AuthRequired())

// 嵌套组
testing := authorized2.Group("testing")
testing.GET( "/analytics" , analyticsEndpoint)
```

#### 单个路由中间件：单个路由中间件仅对一个路由起作用

```go
r := gin.New()
authorized := r.Group("/users")

// 单个路由设置单个中间件
authorized.POST("/login", loginEndpoint)

// 单个路由设置多个中间件
r.GET("/benchmark", MyBenchLogger(), benchEndpoint)
```

#### 总结

上述代码段中，各个中间件作用的路由如下：

- `r.Use(gin.Logger())、r.Use(gin.Recovery())` 将中间件添加在全局路由上，也就是说，所有请求路径以 / 开头的请求都会被 gin.Logger()、gin.Recovery() 中间件按顺序请求；

- AuthRequired() 中间件被添加在了 authorized 路由分组中，也就是所有请求路径以 /users 开头的请求，都会被 AuthRequired() 中间件处理；

- analyticsEndpoint 中间件被添加在了 testing 路由分组中，也就是所有请求路径以 /users/testing 开头的请求，都会被 analyticsEndpoint 中间件处理；

- loginEndpoint 中间件只作用在 POST /users/login 方法。

### 7. Gin中间件开发

**开发一个给请求添加`X-Request-ID`的中间件**

- 在请求中注入 RequestID
- 在日志中打印 RequestID

**查看`miniblog\internal\pkg\middleware\requestid.go`**


**开发一个处理跨域`CORS`的中间件**

- 简单请求:请求方法是 `GET`、`HEAD` 或者 `POST`，并且 HTTP 请求头中只有 `Accept/Accept-Language/Content-Language/Last-Event-ID/Content-Type` 6 种类型，且 `Content-Type` 只能是 `application/x-www-form-urlencoded、multipart/form-data、text、plain` 中的一个值。简单请求会在发送时自动在 HTTP 请求头加上 Origin 字段，来标明当前是哪个源(协议 + 域名 + 端口)，服务端来决定是否放行
- 复杂请求：不是简单请求就是复杂请求，一般`POST`+`json`的都是复杂请求

判断很简单，如果是`OPTION`请求，添加CORS相关的HEADER

**查看`miniblog\internal\pkg\middleware\header.go`**

### 8. 程序优雅关停功能

先来说下，为什么要添加优雅关停能力。在应用程序的生命周期中，新功能发布、Bug 修复、配置变更等，都需要重启服务。在服务进程停止的时候，可能需要做一些处理工作，例如：

1. 正在执行的 HTTP 请求，要等待请求执行完并返回，否则该请求会报错，并产生一些脏数据；

2. 异步处理任务，也需要将缓存中的数据处理完成，否则可能会造成一些数据丢失或者不一致；

3. 关闭数据库连接，否则数据库连接池会保存一个无用的连接，造成宝贵的连接资源浪费。

所以，给应用程序实现优雅关停功能，可以大大提高系统的健壮性。

**核心执行链路**
1. 把启动服务放到 `goroutine` 中
2. 创建 `os.Signal` 类型的 `channel`，用来捕获程序关停信号
3. 调用 `signal.Notify` 函数设置需要捕获的信号，需要设置为 `syscall.SIGINT`, `syscall.SIGTERM` 2 种信号
4. 调用 `<-quit` 阻塞主程序
5. 如果系统收到 `SIGINT` 和 `SIGTERM` 信号，就会往 `quit channel` 中写入一条 `os.Signal` 类型的数据
6. quit 读取到数据，解除阻塞状态
7. 通过 `http.Shutdown` 方法，关停 HTTP 服务

时序图：
```
主 goroutine                    HTTP goroutine
    |                               |
    |--- go ListenAndServe() ------>| (开始监听)
    |                               |
    |<-- 阻塞在 <-quit              | (处理请求中...)
    |                               |
 [Ctrl+C]                           |
    |                               |
    |--- Shutdown(ctx) ------------>| 停止接受新连接
    |                               | 等待处理中的请求完成
    |                               |
    |<-- Shutdown 返回 ------------ |
    |                               |
 log("服务退出")                   goroutine 结束
```

### 9. 优雅处理错误码

**参照[腾讯云API3.0的错误码设计规范](https://github.com/marmotedu/miniblog/blob/master/docs/devel/zh-CN/conversions/error_code.md)，采用二级错误码**

二级错误码
- 语义化： 语义化的错误码，通过错误码名字，就能知道报错的类型
- 更加灵活： 二级错误码格式为 `平台级.资源级`。平台级错误码是固定的，用来指代某一类错误，客户端可使用该错误码，进行通用的错误处理。可使用资源级错误码进行更精准的错误处理。此外，服务端可以根据需要自定义错误码，也可以使用默认的错误码。

这里，我们预定义了以下平台级错误码：

```
错误码	            错误描述	                            错误类型
InternalError	    内部错误	                                1
InvalidParameter	参数错误（包括参数类型、格式、值等错误）   0
AuthFailure	        认证 / 授权错误	                        0
ResourceNotFound	资源不存在	                            0
FailedOperation	    操作失败	                                2
```

**查看`miniblog\internal\pkg\core\core.go`**

![alt text](image-4.png)

## 六、设计业务架构

根据我们所设计的4层架构，我们可知以下依赖关系：Controller 层依赖 Biz 层，Biz 层依赖 Store 层，Store 层依赖数据库，而 Controller层、Biz 层、Store 层都依赖 Model 层

![alt text](image-5.png)

为了能够随时测试我们所开发的代码功能，最好的方式是先开发依赖少的组件，所以开发顺序为：Model 层 -> Store 层 -> Biz 层 -> Controller 层

### 1. model层
- 启动mysql：`net start mysql80`、关闭为 `net stop mysql80`

- 登录mysql：`mysql -h127.0.0.1 -P3306 -uroot -p123456`
![alt text](image-6.png)

- 根据项目的sql生成表：`source configs/miniblog.sql`
    - 把数据库的表生成到项目中 `mysqldump --column-statistics=0 -h127.0.0.1 -uroot --databases miniblog -p123456 --add-drop-database --add-drop-table --add-drop-trigger --add-locks --no-data > configs/miniblog.sql`

- 切换到miniblog数据库：`USE miniblog;`

- 查看miniblog数据库的表列表：`SHOW TABLES;`
![alt text](image-7.png)

- 查看表的内容：`DESC user;` `DESC post;`
![alt text](image-8.png)

- 根据表生成struct：`db2struct --gorm --no-json -H 127.0.0.1 -d miniblog -t user --package model --struct UserM -u root -p '123456' --target=user.go`

### 2. store层

**查看 `miniblog/internal/miniblog/store/store.go`**

- 首先，要创建一个结构体，用来创建 Store 层的实例。自然的，你会想到要在改结构体中包含一个 *gorm.DB 对象，用于与数据库的 CURD；

- 接着，创建一个 New 函数，用来创建 Store 层实例；

- 接着，为了方便直接调用 store 包，引用 Store 层的实例，我们还要设置一个包级别的 Store 实例；

- 最后，为了避免实例被重复创建，通常我们需要使用 sync.Once 来确保实例只被初始化一次。


