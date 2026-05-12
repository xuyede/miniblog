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
