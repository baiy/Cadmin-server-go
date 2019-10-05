package userGroupRelate

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/baiy/Cadmin-service-go/utils/set"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	AdminUserGroupId string `json:"admin_user_group_id"`
	AdminAuthId      string `json:"admin_auth_id"`
}

func AuthIds(userGroupIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_user_group_relate").Select("admin_auth_id").Where(goqu.Ex{
		"admin_user_group_id": userGroupIds,
	}).ScanVals(&ids)
	return ids
}

func UserGroupIds(authIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_user_group_relate").Select("admin_user_group_id").Where(goqu.Ex{
		"admin_auth_id": authIds,
	}).ScanVals(&ids)
	return ids
}

// 用户分组权限检查
func Check(authIds []int, userGroupIds []int) bool {
	if len(authIds) == 0 || len(userGroupIds) == 0 {
		return false
	}
	existAuthIds := AuthIds(userGroupIds)

	if len(existAuthIds) == 0 {
		return false
	}
	return len(set.IntSliceIntersect(existAuthIds, authIds)) != 0
}

func Add(userGroupId, authId int) error {
	_, err := models.Db.Insert("admin_user_group_relate").Rows(
		goqu.Record{"admin_user_group_id": userGroupId, "admin_auth_id": authId},
	).Executor().Exec()
	return err
}

func Remove(userGroupId, authId int) error {
	if userGroupId == 0 && authId == 0 {
		return errors.New("参数错误")
	}
	where := make(goqu.Ex)
	if userGroupId != 0 {
		where["admin_user_group_id"] = userGroupId
	}
	if authId != 0 {
		where["admin_auth_id"] = authId
	}
	_, err := models.Db.Delete("admin_user_group_relate").Where(where).Executor().Exec()
	return err
}
