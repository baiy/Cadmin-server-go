package request

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	Type   string `json:"type"`
	Name   string `json:"name"`
	Action string `json:"action"`
	Call   string `json:"call"`
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
