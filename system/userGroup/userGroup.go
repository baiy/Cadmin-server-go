package userGroup

import (
	"github.com/baiy/Cadmin-server-go/admin"
	"github.com/baiy/Cadmin-server-go/models"
	"github.com/baiy/Cadmin-server-go/models/auth"
	"github.com/baiy/Cadmin-server-go/models/user"
	thisModel "github.com/baiy/Cadmin-server-go/models/userGroup"
	"github.com/baiy/Cadmin-server-go/models/userRelate"
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
		Auth []*auth.Model `db:"-" json:"auth"`
		User []*user.Model `db:"-" json:"user"`
	}, 0)
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["name"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	total, err := param.Select("admin_user_group", &lists, where)
	if err != nil {
		return nil, err
	}

	for index := range lists {
		lists[index].Auth, _ = auth.GetLists(lists[index].AuthIds())
		lists[index].User, _ = user.GetLists(lists[index].UserIds())
	}

	return map[string]interface{}{
		"lists": lists,
		"total": total,
	}, nil
}

func Save(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id          int    `form:"id"`
		Name        string `form:"name" validate:"required"`
		Description string `form:"description"`
	})

	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	if param.Id == 0 {
		return nil, thisModel.Add(param.Name, param.Description)
	}
	return nil, thisModel.Updata(param.Id, param.Name, param.Description)
}

func Remove(context *admin.Context) (interface{}, error) {
	id, err := context.InputInt("id")
	if err != nil {
		return nil, err
	}
	return nil, thisModel.Remove(id)
}

func GetUser(context *admin.Context) (interface{}, error) {
	param := new(struct {
		utils.Page
		Id      int    `form:"id" validate:"required"`
		Keyword string `form:"keyword"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	current, err := thisModel.GetById(param.Id)
	if err != nil {
		return nil, err
	}

	noAssign := make([]user.Model, 0)
	exist := current.UserIds()
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["username"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	if len(exist) != 0 {
		where["id"] = goqu.Op{"notin": exist}
	}
	total, err := param.Select("admin_user", &noAssign, where)
	if err != nil {
		return nil, err
	}
	assign := make([]user.Model, 0)
	err = models.Db.From(goqu.T("admin_user").As("user")).
		Select("user.*").
		InnerJoin(
			goqu.T("admin_user_relate").As("rel"),
			goqu.On(goqu.Ex{
				"rel.admin_user_id": goqu.I("user.id"),
			}),
		).
		Where(goqu.Ex{"rel.admin_user_group_id": param.Id}).
		Order(goqu.I("rel.id").Desc()).ScanStructs(&assign)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"lists": map[string]interface{}{
			"assign":   assign,
			"noAssign": noAssign,
		},
		"total": total,
	}, nil
}

func AssignUser(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id     int `form:"id" validate:"required"`
		UserId int `form:"userId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, userRelate.Add(param.Id, param.UserId)
}

func RemoveUser(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id     int `form:"id" validate:"required"`
		UserId int `form:"userId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, userRelate.Remove(param.Id, param.UserId)
}
