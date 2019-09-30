package auth

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/baiy/Cadmin-service-go/models/menuRelate"
	"github.com/baiy/Cadmin-service-go/models/requestRelate"
	"github.com/baiy/Cadmin-service-go/models/userGroupRelate"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (m Model) RequestIds() []int {
	return requestRelate.RequestIds([]int{m.Id})
}

func (m Model) MenuIds() []int {
	return menuRelate.MenuIds([]int{m.Id})
}

func (m Model) UserGroupIds() []int {
	return userGroupRelate.UserGroupIds([]int{m.Id})
}

func Add(name, description string) error {
	_, err := models.Db.Insert("admin_auth").Rows(
		goqu.Record{"name": name, "description": description},
	).Executor().Exec()
	return err
}

func Updata(id int, name, description string) error {
	_, err := models.Db.Update("admin_auth").Set(goqu.Record{"name": name, "description": description}).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	return err
}

func Remove(id int) error {
	_, err := models.Db.Delete("admin_auth").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err == nil {
		_ = requestRelate.Remove(0, id)
		_ = userGroupRelate.Remove(0, id)
		_ = menuRelate.Remove(0, id)
	}
	return err
}

func GetById(id int) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_auth").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("权限不存在")
		}
	}
	return
}

func GetLists(ids []int) ([]*Model, error) {
	model := make([]*Model, 0)
	if len(ids) == 0 {
		return model, nil
	}
	err := models.Db.From("admin_auth").Where(goqu.Ex{
		"id": ids,
	}).ScanStructs(&model)
	return model, err
}
