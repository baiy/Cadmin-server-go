package models

import (
	"database/sql"
	"github.com/baiy/Cadmin-server-go/models/utils"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

var Db *goqu.Database

type Model struct {
	Id         int        `db:"id" json:"id" goqu:"skipinsert,skipupdate"`
	CreateTime utils.Time `db:"create_time" json:"create_time" goqu:"skipinsert,skipupdate"`
	UpdateTime utils.Time `db:"update_time" json:"update_time" goqu:"skipinsert,skipupdate"`
}

func InitDb(d *sql.DB) {
	// 设置列名自动转换规则
	goqu.SetColumnRenameFunction(func(s string) string {
		snake := regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}")
		return strings.ToLower(snake)
	})
	dialect := goqu.Dialect("mysql")
	Db = dialect.DB(d)
}