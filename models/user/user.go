package user

import (
	"errors"
	"fmt"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/baiy/Cadmin-service-go/models/menuRelate"
	"github.com/baiy/Cadmin-service-go/models/requestRelate"
	"github.com/baiy/Cadmin-service-go/models/userGroupRelate"
	"github.com/baiy/Cadmin-service-go/models/userRelate"
	"github.com/baiy/Cadmin-service-go/models/utils"
	"github.com/doug-martin/goqu/v9"
	"time"
)

const (
	Enable  = 1 // 启用
	Disable = 2 // 禁用
)

type Model struct {
	models.Model
	Username      string      `json:"username"`
	Password      string      `json:"-"`
	LastLoginIp   string      `json:"last_login_ip"`
	LastLoginTime *utils.Time `json:"last_login_time"`
	Status        int         `json:"status"`
}

func (m *Model) IsDisabled() bool {
	return m.Status == Disable
}

func (m *Model) LoginUpdate(ip string) {
	_, _ = models.Db.Update("admin_user").Where(goqu.Ex{
		"id": m.Id,
	}).Set(map[string]interface{}{
		"last_login_ip":   ip,
		"last_login_time": time.Now().Format("2006-01-02 15:04:05"),
	}).Executor().Exec()
}

func (m Model) UserGroupIds() []int {
	return userRelate.GroupIds([]int{m.Id})
}

func (m Model) AuthIds() []int {
	return userGroupRelate.AuthIds(m.UserGroupIds())
}

func (m Model) MenuIds() []int {
	return menuRelate.MenuIds(m.AuthIds())
}

func (m Model) RequestIds() []int {
	return menuRelate.MenuIds(m.AuthIds())
}

func Add(username, password string, status int) error {
	exist, _ := GetByUserName(username)
	if exist.Id > 0 {
		return errors.New(fmt.Sprintf("[%s] 用户已经存在", username))
	}
	_, err := models.Db.Insert("admin_user").Rows(
		goqu.Record{"username": username, "password": password, "status": status},
	).Executor().Exec()
	return err
}

func Updata(id int, username, password string, status int) error {
	exist, _ := GetByUserName(username)
	if exist.Id > 0 && exist.Id != id {
		return errors.New(fmt.Sprintf("[%s] 用户已经存在", username))
	}
	record := goqu.Record{"username": username, "status": status}
	if password != "" {
		record["password"] = password
	}
	_, err := models.Db.Update("admin_user").Set(record).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	return err
}

func Remove(id int) error {
	_, err := models.Db.Delete("admin_user").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err == nil {
		err = userRelate.Remove(0, id)
	}
	return err
}

// 检查用户权限
func CheckAuth(id int, requestId int) error {
	userGroupIds := userRelate.GroupIds([]int{id})
	if len(userGroupIds) == 0 {
		return errors.New("用户未分配用户组")
	}

	authIds := requestRelate.AuthIds([]int{requestId})
	if len(authIds) == 0 {
		return errors.New("请求未分配权限组")
	}

	if !userGroupRelate.Check(authIds, userGroupIds) {
		return errors.New("暂无权限")
	}
	return nil
}

func All() ([]*Model, error) {
	m := make([]*Model, 0)
	err := models.Db.From("admin_user").ScanStructs(&m)
	return m, err
}

func GetById(id int) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_user").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("用户不存在")
		}
	}
	return
}

func GetByUserName(username string) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_user").Where(goqu.Ex{
		"username": username,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("用户不存在")
		}
	}
	return
}

func GetLists(ids []int) ([]*Model, error) {
	model := make([]*Model, 0)
	if len(ids) == 0 {
		return model, nil
	}
	err := models.Db.From("admin_user").Where(goqu.Ex{
		"id": ids,
	}).ScanStructs(&model)
	return model, err
}
