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