package request

import (
	"github.com/baiy/Cadmin-server-go/admin"
	"github.com/baiy/Cadmin-server-go/models/auth"
	thisModel "github.com/baiy/Cadmin-server-go/models/request"
	"github.com/baiy/Cadmin-server-go/system/utils"
	"github.com/doug-martin/goqu/v9"
)

type DispatchItem struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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
	}, 0)
	where := make(goqu.ExOr)
	if param.Keyword != "" {
		where["name"] = goqu.Op{"like": "%" + param.Keyword + "%"}
		where["action"] = goqu.Op{"like": "%" + param.Keyword + "%"}
		where["call"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	total, err := param.Select("admin_request", &lists, where)
	if err != nil {
		return nil, err
	}

	for index := range lists {
		lists[index].Auth, _ = auth.GetLists(lists[index].AuthIds())
	}

	return map[string]interface{}{
		"lists": lists,
		"total": total,
	}, nil
}

func Save(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id     int    `form:"id"`
		Name   string `form:"name" validate:"required"`
		Action string `form:"action" validate:"required"`
		Type   string `form:"type" validate:"required"`
		Call   string `form:"call"`
	})

	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	if param.Id == 0 {
		return nil, thisModel.Add(param.Name, param.Action, param.Type, param.Call)
	}
	return nil, thisModel.Updata(param.Id, param.Name, param.Action, param.Type, param.Call)
}

func Remove(context *admin.Context) (interface{}, error) {
	id, err := context.InputInt("id")
	if err != nil {
		return nil, err
	}
	return nil, thisModel.Remove(id)
}
func Type(context *admin.Context) (interface{}, error) {
	lists := make([]DispatchItem, admin.AllDispatcherLength())
	i := 0
	for type_, value := range admin.AllDispatcher() {
		lists[i] = DispatchItem{
			Type:        type_,
			Name:        value.Name(),
			Description: value.Description(),
		}
		i++
	}
	return lists, nil
}
