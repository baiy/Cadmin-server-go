Cadmin golang 服务端 

> 项目地址: [https://github.com/baiy/Cadmin-server-go](https://github.com/baiy/Cadmin-server-go)

### 特点

1. 为便于给现有系统加入后台管理功能和加快新系统开发, 后台核心系统尽可能的减少依赖, 不侵入外层业务系统.
2. 对请求处理按照请求类型可自定义`请求调度类`,便于不用业务系统使用和开发. 

### 安装
```
go get -u github.com/baiy/Cadmin-server-go
```

### 数据库

详见 [数据库结构](server/db.md) 一章

### 使用方法
> 在代码安装和数据库导入完毕后, 接下来需要将后台系统的入口代码嵌入当前系统的合适位置, 并进行相应的配置

#### 入口代码示例
[example.go](https://raw.githubusercontent.com/baiy/Cadmin-server-go/master/example.go ':include :type=go')

### 自定义用户密码生成策略

1. 实现 `github.com/baiy/Cadmin-service-go/admin.Passwrod` 接口
2. 注册密码生成器,使用`github.com/baiy/Cadmin-service-go/admin.RegisterPassword()`

系统内置密码生成器: <https://github.com/baiy/Cadmin-server-go/blob/master/admin/password.go>

> 内置密码生成规则: `base64_encode(hash('sha256',hash("sha256", $password.$salt,FALSE).$salt,FALSE).'|'.$salt);`

### 请求调度器开发

1. 实现 `github.com/baiy/Cadmin-service-go/admin.Dispatch` 接口
2. 注册调度器,使用`github.com/baiy/Cadmin-service-go/admin.RegisterDispatch()`

内置调度器: <https://github.com/baiy/Cadmin-server-go/blob/master/admin/dispatch.go>

### 使用默认调度器开发

```go
package router

import "github.com/baiy/Cadmin-service-go/admin"
admin.RegisterDefaultDispatcherHandleMethod(map[string]admin.DefaultDispatcherHandleMethod{
    'request.call':func (context *admin.Context) (interface{}, error) {return nil,nil}
})
```

> `request.call` 对应后台请求中的类型配置值