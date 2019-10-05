package main

import (
	"database/sql"
	"fmt"
	"github.com/baiy/Cadmin-service-go/admin"
	_ "github.com/baiy/Cadmin-service-go/system" // 加载Cadmin内置action路由
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var db *sql.DB

func initDb() {
	var err error
	db, err = sql.Open("mysql", "root:root@/admin_api_new?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}

func main() {
	initDb()
	admin.InjectDb(db)
	http.HandleFunc("/admin/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		context := admin.NewContext(writer, request)
		if err := context.Output(); err != nil {
			fmt.Println(err)
		}
	})
	_ = http.ListenAndServe("127.0.0.1:8001", nil)
}
