package user

import (
	"errors"
	"github.com/baiy/Cadmin-server-go/admin"
	thisModel "github.com/baiy/Cadmin-server-go/models/user"
	"github.com/baiy/Cadmin-server-go/models/userGroup"
	"github.com/baiy/Cadmin-server-go/system/utils"
	"github.com/doug-martin/goqu/v9"
)

func Lists(context *admin.Context) (interface{}, error) {
	param := new(struct {
		utils.Page
		Keyword string `form:"keyword"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	lists := make([]struct {
		thisModel.Model
		UserGroup []*userGroup.Model `db:"-" json:"userGroup"`
	}, 0)
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["username"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	total, err := param.Select("admin_user", &lists, where)
	if err != nil {
		return nil, err
	}

	for index := range lists {
		lists[index].UserGroup, _ = userGroup.GetLists(lists[index].UserGroupIds())
	}

	return map[string]interface{}{
		"lists": lists,
		"total": total,
	}, nil
}

func Save(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id          int    `form:"id"`
		Username    string `form:"username" validate:"required"`
		Password    string `form:"password"`
		Description string `form:"description"`
		Status      int    `form:"status"  validate:"required"`
	})

	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	password := ""
	if param.Password != "" {
		password = string(admin.Passworder.Hash([]byte(param.Password)))
	}
	if param.Id == 0 {
		if param.Password == "" {
			return nil, errors.New("添加用户密码不能为空")
		}
		return nil, thisModel.Add(param.Username, password, param.Status,param.Description)
	}
	return nil, thisModel.Updata(param.Id, param.Username, password, param.Status,param.Description)
}

func Remove(context *admin.Context) (interface{}, error) {
	id, err := context.InputInt("id")
	if err != nil {
		return nil, err
	}
	return nil, thisModel.Remove(id)
}
