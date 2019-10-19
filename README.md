Cadmin golang 服务端 

> 项目地址: [[github](https://github.com/baiy/Cadmin-server-go)] [[gitee](https://gitee.com/baiy/Cadmin-server-go)]
>
> 在线文档地址: <https://baiy.github.io/Cadmin/>

### 特点

1. 为便于给现有系统加入后台管理功能和加快新系统开发, 后台核心系统尽可能的减少依赖, 不侵入外层业务系统.
2. 对请求处理按照请求类型可自定义`请求调度类`,便于不用业务系统使用和开发. 

### 安装

```
go get -u github.com/baiy/Cadmin-server-go
```

### 数据库

详见 [数据库结构](https://baiy.github.io/Cadmin/#/server/db) 一章

### 使用方法
> 在代码安装和数据库导入完毕后, 接下来需要将后台系统的入口代码嵌入当前系统的合适位置, 并进行相应的配置

#### 原生

```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/baiy/Cadmin-server-go/admin"
	_ "github.com/baiy/Cadmin-server-go/system" // 注册内置请求处理方法
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var db *sql.DB


// 初始化数据库
func initDb() {
	var err error
	db, err = sql.Open("mysql", "root:root@/admin_api_new?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}

func main() {
	initDb()

	// 设置数据库操作对象
	admin.SetDb(db)
	// [可选] 注册自定义调度器
	// admin.RegisterDispatch()
	// [可选] 设置自定义密码生成器
	// admin.RegisterPassword()
	// [可选] 无需校验权限的api
	// admin.AddNoCheckLoginRequestId()
	// [可选] 只需登录即可访问的api
	// admin.AddOnlyLoginRequestId()
	// [可选] 设置请求标识变量名
	// admin.ActionName = ""
	// [可选] 设置请求token变量名
	// admin.TokenName = ""
	// [可选] 设置请求日志记录回调函数
	//admin.RegisterLogCallback()

	http.HandleFunc("/api/admin/", func(writer http.ResponseWriter, request *http.Request) {
		// 前后端分离项目一般会有跨域问题 自行处理
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		context := admin.NewContext(writer, request)
		if err := context.Output(); err != nil {
			fmt.Println(err)
		}
	})
	_ = http.ListenAndServe("127.0.0.1:8001", nil)
}

```

#### Beego

```go
package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/baiy/Cadmin-server-go/admin"
	_ "github.com/baiy/Cadmin-server-go/system" // 注册内置请求处理方法
	_ "github.com/go-sql-driver/mysql"          // import your used driver
)

func main() {
	orm.RegisterDataBase("default", "mysql", "root:root@/admin_api_new?charset=utf8mb4&parseTime=True&loc=Local", 30)
	db, _ := orm.GetDB("default")
	admin.SetDb(db)
	// 省略其他配置代码 查看原生示例代码

	beego.Any("/api/admin/", func(ctx *context.Context) {
		ctx.ResponseWriter.ResponseWriter.Header().Add("Access-Control-Allow-Origin", "*")
		ctx.ResponseWriter.ResponseWriter.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		c := admin.NewContext(ctx.ResponseWriter.ResponseWriter, ctx.Request)
		if err := c.Output(); err != nil {
			fmt.Println(err)
		}
	})
	beego.Run("127.0.0.1:8001")
}
```

#### Gorm

```go
package main

import (
	"fmt"
	"github.com/baiy/Cadmin-server-go/admin"
	_ "github.com/baiy/Cadmin-server-go/system" // 注册内置请求处理方法
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
)

func main() {
	db, _ := gorm.Open("mysql", "root:root@/admin_api_new?charset=utf8mb4&parseTime=True&loc=Local")
	defer db.Close()

	admin.SetDb(db.DB())
	// 省略其他配置代码 查看原生示例代码
	http.HandleFunc("/api/admin/", func(writer http.ResponseWriter, request *http.Request) {
		// 前后端分离项目一般会有跨域问题 自行处理
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		context := admin.NewContext(writer, request)
		if err := context.Output(); err != nil {
			fmt.Println(err)
		}
	})
	_ = http.ListenAndServe("127.0.0.1:8001", nil)

}
```

#### Xorm

```go
package main

import (
	"fmt"
	"github.com/baiy/Cadmin-server-go/admin"
	_ "github.com/baiy/Cadmin-server-go/system" // 注册内置请求处理方法
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/http"
)

func main() {
	engine, _ := xorm.NewEngine("mysql", "root:root@/admin_api_new?charset=utf8mb4&parseTime=True&loc=Local")
	admin.SetDb(engine.DB().DB)
	// 省略其他配置代码 查看原生示例代码
	http.HandleFunc("/api/admin/", func(writer http.ResponseWriter, request *http.Request) {
		// 前后端分离项目一般会有跨域问题 自行处理
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		context := admin.NewContext(writer, request)
		if err := context.Output(); err != nil {
			fmt.Println(err)
		}
	})
	_ = http.ListenAndServe("127.0.0.1:8001", nil)
}
```


### 自定义用户密码生成策略

1. 实现 `github.com/baiy/Cadmin-service-go/admin.Passwrod` 接口
2. 注册密码生成器,使用`github.com/baiy/Cadmin-service-go/admin.RegisterPassword()`

系统内置密码生成器: <https://github.com/baiy/Cadmin-server-go/blob/master/admin/password.go>

### 请求调度器开发

1. 实现 `github.com/baiy/Cadmin-service-go/admin.Dispatch` 接口
2. 注册调度器,使用`github.com/baiy/Cadmin-service-go/admin.RegisterDispatch()`

内置调度器: <https://github.com/baiy/Cadmin-server-go/blob/master/admin/dispatch.go>

### 使用默认调度器开发

```go
package router

import (
    "github.com/baiy/Cadmin-service-go/admin"
)

func init(){
    admin.RegisterDefaultDispatcherHandleMethod(map[string]admin.DefaultDispatcherHandleMethod{
        "request.call":func (context *admin.Context) (interface{}, error) {return nil,nil},
    })
}
```

> `request.call` 对应后台请求中的类型配置值