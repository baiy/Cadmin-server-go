package auth

import (
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	Name        string `json:"name"`
	Description string `json:"description"`
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
