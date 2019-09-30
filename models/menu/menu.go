package menu

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/baiy/Cadmin-service-go/models/menuRelate"
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

func Add(parentId int, name, url, icon, description string) error {
	if parentId != 0 {
		exist, err := GetById(parentId)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		if exist.Url != "" {
			return errors.New("父菜单不是目录类型菜单")
		}
	}
	_, err := models.Db.Insert("admin_menu").Rows(
		goqu.Record{
			"parent_id":   parentId,
			"name":        name,
			"url":         url,
			"icon":        icon,
			"description": description,
			"sort":        getNewSort(parentId),
		},
	).Executor().Exec()
	return err
}

func Updata(id int, parentId int, name, url, icon, description string) error {
	if parentId != 0 {
		exist, err := GetById(parentId)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		if exist.Url != "" {
			return errors.New("父菜单不是目录类型菜单")
		}
	}
	_, err := models.Db.Update("admin_menu").Where(goqu.Ex{"id": id}).Set(
		goqu.Record{
			"parent_id":   parentId,
			"name":        name,
			"url":         url,
			"icon":        icon,
			"description": description,
		},
	).Executor().Exec()
	return err
}

func Remove(id int) error {
	_, err := models.Db.Delete("admin_menu").Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	if err == nil {
		_ = menuRelate.Remove(id, 0)
	}
	return err
}

func All() ([]*Model, error) {
	m := make([]*Model, 0)
	err := models.Db.From("admin_menu").ScanStructs(&m)
	return m, err
}

func Sort(id, sort int) error {
	_, err := models.Db.Update("admin_menu").Set(goqu.Record{"sort": sort}).Where(goqu.Ex{
		"id": id,
	}).Executor().Exec()
	return err
}

func getNewSort(parentId int) int {
	model := new(Model)
	found, err := models.Db.From("admin_menu").Where(goqu.Ex{
		"parent_id": parentId,
	}).Order(goqu.I("sort").Desc()).ScanStruct(model)
	if err == nil && found {
		return model.Sort + 1
	}
	return 0
}

func GetById(id int) (model *Model, err error) {
	model = new(Model)
	found, err := models.Db.From("admin_menu").Where(goqu.Ex{
		"id": id,
	}).ScanStruct(model)
	if err == nil {
		if !found {
			err = errors.New("菜单不存在")
		}
	}
	return
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
