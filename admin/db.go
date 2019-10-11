package admin

import (
	"database/sql"
	"github.com/baiy/Cadmin-server-go/models"
)

var Db *sql.DB

// 设置数据库操作对象
func SetDb(d *sql.DB) {
	Db = d
	models.InitDb(Db)
}
