package menu

import (
	"github.com/baiy/Cadmin-server-go/admin"
	thisModel "github.com/baiy/Cadmin-server-go/models/menu"
)

func Lists(context *admin.Context) (interface{}, error) {
	return thisModel.All()
}

func Save(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id          int    `form:"id"`
		ParentId    int    `form:"parent_id" validate:"min=0"`
		Name        string `form:"name" validate:"required"`
		Url         string `form:"url"`
		Icon        string `form:"icon"`
		Description string `form:"description"`
	})

	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	if param.Id == 0 {
		return nil, thisModel.Add(param.ParentId, param.Name, param.Url, param.Icon, param.Description)
	}
	return nil, thisModel.Updata(param.Id, param.ParentId, param.Name, param.Url, param.Icon, param.Description)
}

func Remove(context *admin.Context) (interface{}, error) {
	id, err := context.InputInt("id")
	if err != nil {
		return nil, err
	}
	return nil, thisModel.Remove(id)
}

func Sort(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Menus map[int]map[string]int `form:"menus"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	for _, item := range param.Menus {
		err = thisModel.Sort(item["id"], item["sort"])
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
