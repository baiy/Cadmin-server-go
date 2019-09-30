package requestRelate

import (
	"errors"
	"github.com/baiy/Cadmin-service-go/models"
	"github.com/doug-martin/goqu/v9"
)

type Model struct {
	models.Model
	AdminRequestId int `json:"admin_request_id"`
	AdminAuthId    int `json:"admin_auth_id"`
}

func AuthIds(requestIds []int) []int {
	ids := make([]int, 0)
	_ = models.Db.From("admin_request_relate").Select("admin_auth_id").Where(goqu.Ex{
		"admin_request_id": requestIds,
	}).ScanVals(&ids)
	return ids
}

func Remove(requestId, authId int) error {
	if requestId == 0 && authId == 0 {
		return errors.New("参数错误")
	}
	where := make(goqu.Ex)
	if requestId != 0 {
		where["admin_request_id"] = requestId
	}
	if authId != 0 {
		where["admin_auth_id"] = authId
	}
	_, err := models.Db.Delete("admin_request_relate").Where(where).Executor().Exec()
	return err
}
