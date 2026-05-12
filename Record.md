### 一、初始化项目

#### 1. 初始化GO模块
- go mod init <project>

#### 2. 初始化Git
- 使用 [.gitignore](https://www.toptal.com/developers/gitignore)快速生成相关项目的模版

#### 3. 热加载工具

- 使用 [air](https://github.com/air-verse/air/blob/master/README-zh_cn.md)

```go install github.com/cosmtrek/air@latest```

**配置 .air.toml 文件**

```执行：air -c .air.toml```

#### 4. 使用Swagger编写接口文档
1. ```go install github.com/go-swagger/go-swagger/cmd/swagger@latest```

2. ```swagger serve -F=swagger --no-open --port 65534 ./api/openapi/openapi.yaml```

3. 在浏览器中打开 http://localhost:65534/docs

#### 5. 添加 LICENSE
1. ```go install github.com/nishanths/license/v5@latest```

2. ```license -n 'Deye(许业德) <xuyede77@gmail.com>' -o LICENSE mit```

#### 6. 给源文件添加版本声明

1. ```go install github.com/marmotedu/addlicense@latest```

2. ```addlicense -v -f ./scripts/boilerplate.txt --skip-dirs=third_party,vendor,_output .cmd/miniblog/main.go```

#### 7. 使用Makefile实现上面的步骤
**配置 Makefile 文件**

- Makefile [学习](https://github.com/marmotedu/geekbang-go/blob/master/makefile/Makefile%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86.md)

### 二、应用组成及构建

#### 1. 引用配置

- 命令行选项、命令行参数： 选择 [pflag](https://github.com/spf13/pflag)；

- 配置文件： 建议选择支持多种配置文件格式的包，也即从：viper、configor、koanf、config 中选择其一，毫无疑问 [viper](https://github.com/spf13/viper) 胜出：
  - viper 有 21000+ 的 Star 数，比其它包更受欢迎；
  - viper 功能强大，并且经过很多大型项目验证过；
  - viper 同时也可以实现分布式配置中心的功能。

- 环境变量： 如果环境变量不多，可以使用 os.Getenv，如果环境变量很多，可以使用 [envconfig](https://github.com/kelseyhightower/envconfig) 直接将环境变量读取到 Go 结构体变量中。

- 配置中心： 可根据需要选择 viper、apollo、[etcd](https://pkg.go.dev/go.etcd.io/etcd/client/v3)、[consul](https://github.com/hashicorp/consul) 等。其实一般的项目不需要引入配置中心，因为使用配置中心，会带来一些部署、维护的复杂度。

**结论：直接使用`pflag` + `viper`**

#### 2. 应用业务逻辑

应用的业务逻辑根据业务的不同差别很大。一般而言，一个 Go 应用中会执行以下类别的业务逻辑处理（可能会用到其中一个或多个）：

- 初始化缓存；

- 初始化并创建各类数据库客户端，例如：Redis、MySQL、Kafka、MongoDB、Etcd 等；

- 初始化并创建其他服务的客户端等；

- 初始化并启动Web服务，例如：HTTP、HTTPS、GRPC；

- 启动异步任务，这些异步任务可以执行任何业务需要的操作，例如：watch kube-apiserver、定期从第三方服务拉取数据，并缓存、注册 /metrics 并监听指定的端口、启动 kafka 消费队列等等；

- 执行特定的业务处理，并退出程序；

#### 3. 应用启动框架

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

#### 4. 最佳构建方法

**结论: 使用 `pflag`、`viper`、`cobra` 来构建一个强大的应用程序**

- pflag - 给应用添加命令行标识 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8Pflag%E7%BB%99%E5%BA%94%E7%94%A8%E6%B7%BB%E5%8A%A0%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A0%87%E8%AF%86.md)
- viper - 配置解析 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E9%85%8D%E7%BD%AE%E8%A7%A3%E6%9E%90%E7%A5%9E%E5%99%A8-Viper%E5%85%A8%E8%A7%A3.md)
- cobra - 命令行框架 [学习](https://github.com/marmotedu/geekbang-go/blob/master/%E7%8E%B0%E4%BB%A3%E5%8C%96%E7%9A%84%E5%91%BD%E4%BB%A4%E8%A1%8C%E6%A1%86%E6%9E%B6-Cobra%E5%85%A8%E8%A7%A3.md)
