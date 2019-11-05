package userRelate

import (
	"errors"
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	AdminUserId      int `json:"admin_user_id"`
	AdminUserGroupId int `json:"admin_user_group_id"`
}

func GroupIds(userIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_user_relate").Select("admin_user_group_id").Where(goqu.Ex{
		"admin_user_id": userIds,
	}).ScanVals(&ids)
	return ids
}

func UserIds(groupIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_user_relate").Select("admin_user_id").Where(goqu.Ex{
		"admin_user_group_id": groupIds,
	}).ScanVals(&ids)
	return ids
}

func Add(groupId, userId int) error {
	_, err := models.Db.Insert("admin_user_relate").Rows(
		goqu.Record{"admin_user_id": userId, "admin_user_group_id": groupId},
	).Executor().Exec()
	return err
}

func Remove(groupId, userId int) error {
	if groupId == 0 && userId == 0 {
		return errors.New("参数错误")
	}
	where := make(goqu.Ex)
	if groupId != 0 {
		where["admin_user_group_id"] = groupId
	}
	if userId != 0 {
		where["admin_user_id"] = userId
	}
	_, err := models.Db.Delete("admin_user_relate").Where(where).Executor().Exec()
	return err
}
