package userGroup

import (
	"errors"
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/baiy/Cadmin-server-go/models/userGroupRelate"
	"github.com/baiy/Cadmin-server-go/models/userRelate"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m Model) AuthIds() []int {
	return userGroupRelate.AuthIds([]int{m.Id})
}

func (m Model) UserIds() []int {
	return userRelate.UserIds([]int{m.Id})
}

func Add(name, description string) error {
	_, err := models.Db.Insert("admin_user_group").Rows(
		goqu.Record{"name": name, "description": description},
	).Executor().Exec()
	return err
}

func Updata(id int, name, description string) error {
	_, err := models.Db.Update("admin_user_group").Set(goqu.Record{"name": name, "description": description}).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	return err
}

func Remove(id int) error {
	_, err := models.Db.Delete("admin_user_group").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err == nil {
		_ = userRelate.Remove(id, 0)
		_ = userGroupRelate.Remove(id, 0)
	}
	return err
}

func GetById(id int) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_user_group").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("用户组不存在")
		}
	}
	return
}

func GetLists(ids []int) ([]*Model, error) {
	model := make([]*Model, 0)
	if len(ids) == 0 {
		return model, nil
	}
	err := models.Db.From("admin_user_group").Where(goqu.Ex{
		"id": ids,
	}).ScanStructs(&model)
	return model, err
}
