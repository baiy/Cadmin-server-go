package menu

import (
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	ParentId    int    `json:"parent_id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

func GetLists(ids []int) ([]*Model, error) {
	model := make([]*Model, 0)
	if len(ids) == 0 {
		return model, nil
	}
	err := models.Db.From("admin_menu").Where(goqu.Ex{
		"id": ids,
	}).ScanStructs(&model)
	return model, err
}
