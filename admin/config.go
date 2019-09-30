package admin

import (
	"database/sql"
	"github.com/baiy/Cadmin-service-go/models"
)

var noCheckLoginRequestIds = []int{1}
var onlyLoginRequestIds = []int{2, 3}

var Config = &struct {
	Db                     *sql.DB
	ActionName             string
	TokenName              string
	NoCheckLoginRequestIds []int                    // 无需登录的请求ID
	OnlyLoginRequestIds    []int                    // 无需检查权限/只需要登录的请求ID
	LogCallback            func(content LogContent) // 日志写入回调函数
	PasswodHash            Passwrod                 // 密码生成对象
}{
	ActionName:             "_action",
	TokenName:              "_token",
	NoCheckLoginRequestIds: noCheckLoginRequestIds,
	OnlyLoginRequestIds:    onlyLoginRequestIds,
	PasswodHash:            &PasswrodDefault{},
}

func InjectDb(d *sql.DB) {
	Config.Db = d
	models.InitDb(Config.Db)
}

// 添加无需登录的请求ID
func AddNoCheckLoginRequestId(ids ...int) {
	Config.NoCheckLoginRequestIds = append(Config.NoCheckLoginRequestIds, ids...)
}

// 无需检查权限/只需要登录的请求ID
func AddOnlyLoginRequestId(ids ...int) {
	Config.OnlyLoginRequestIds = append(Config.OnlyLoginRequestIds, ids...)
}
