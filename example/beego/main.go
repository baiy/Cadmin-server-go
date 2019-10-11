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