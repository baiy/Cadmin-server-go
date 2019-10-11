package menuRelate

import (
	"errors"
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	AdminMenuId int `json:"admin_menu_id"`
	AdminAuthId int `json:"admin_auth_id"`
}

func MenuIds(authIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_menu_relate").Select("admin_menu_id").Where(goqu.Ex{
		"admin_auth_id": authIds,
	}).ScanVals(&ids)
	return ids
}

func Add(menuId, authId int) error {
	_, err := models.Db.Insert("admin_menu_relate").Rows(
		goqu.Record{"admin_menu_id": menuId, "admin_auth_id": authId},
	).Executor().Exec()
	return err
}

func AddMultiple(menuIds []int, authId int) error {
	rows := make([]goqu.Record, len(menuIds))
	for key, menuId := range menuIds {
		rows[key] = goqu.Record{
			"admin_menu_id": menuId,
			"admin_auth_id": authId,
		}
	}
	_, err := models.Db.Insert("admin_menu_relate").Rows(rows).Executor().Exec()
	return err
}

func RemoveMultiple(menuIds []int, authId int) error {
	_, err := models.Db.Delete("admin_menu_relate").Where(goqu.Ex{
		"admin_auth_id": authId,
		"admin_menu_id": menuIds,
	}).Executor().Exec()
	return err
}

func Remove(menuId, authId int) error {
	if menuId == 0 && authId == 0 {
		return errors.New("参数错误")
	}
	where := make(goqu.Ex)
	if menuId != 0 {
		where["admin_menu_id"] = menuId
	}
	if authId != 0 {
		where["admin_auth_id"] = authId
	}
	_, err := models.Db.Delete("admin_menu_relate").Where(where).Executor().Exec()
	return err
}
