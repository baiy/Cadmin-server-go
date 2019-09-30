package request

import (
	"errors"
	"fmt"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/baiy/Cadmin-service-go/models/requestRelate"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	Type   string `json:"type"`
	Name   string `json:"name"`
	Action string `json:"action"`
	Call   string `json:"call"`
}

func (m Model) AuthIds() []int {
	return requestRelate.AuthIds([]int{m.Id})
}

func Add(name, action, type_, call string) error {
	exist, _ := GetByAction(action)
	if exist.Id > 0 {
		return errors.New(fmt.Sprintf("[%s] 请求已经存在", action))
	}
	_, err := models.Db.Insert("admin_request").Rows(
		goqu.Record{"name": name, "action": action, "type": type_, "call": call},
	).Executor().Exec()
	return err
}

func Updata(id int, name, action, type_, call string) error {
	exist, _ := GetByAction(action)
	if exist.Id > 0 && exist.Id != id {
		return errors.New(fmt.Sprintf("[%s] 请求已经存在", action))
	}
	_, err := models.Db.Update("admin_request").Where(goqu.Ex{"id": id}).Set(
		goqu.Record{"name": name, "action": action, "type": type_, "call": call},
	).Executor().Exec()
	return err
}

func Remove(id int) error {
	_, err := models.Db.Delete("admin_request").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err == nil {
		_ = requestRelate.Remove(id, 0)
	}
	return err
}

func GetByAction(action string) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_request").Where(goqu.Ex{
		"action": action,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("请求不存在")
		}
	}
	return
}

func GetLists(ids []int) ([]*Model, error) {
	model := make([]*Model, 0)
	if len(ids) == 0 {
		return model, nil
	}
	if len(ids) == 0 {
		return model, nil
	}
	err := models.Db.From("admin_request").Where(goqu.Ex{
		"id": ids,
	}).ScanStructs(&model)
	return model, err
}
