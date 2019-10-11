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
