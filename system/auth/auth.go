package auth

import (
	"github.com/baiy/Cadmin-service-go/admin"
	"github.com/baiy/Cadmin-service-go/models"
	thisModel "github.com/baiy/Cadmin-service-go/models/auth"
	"github.com/baiy/Cadmin-service-go/models/menu"
	"github.com/baiy/Cadmin-service-go/models/menuRelate"
	"github.com/baiy/Cadmin-service-go/models/request"
	"github.com/baiy/Cadmin-service-go/models/requestRelate"
	"github.com/baiy/Cadmin-service-go/models/userGroup"
	"github.com/baiy/Cadmin-service-go/models/userGroupRelate"
	"github.com/baiy/Cadmin-service-go/system/utils"
	"github.com/baiy/Cadmin-service-go/utils/set"
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
		Request   []*request.Model   `db:"-" json:"request"`
		Menu      []*menu.Model      `db:"-" json:"menu"`
		UserGroup []*userGroup.Model `db:"-" json:"userGroup"`
	}, 0)
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["name"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	total, err := param.Select("admin_auth", &lists, where)
	if err != nil {
		return nil, err
	}

	for index := range lists {
		lists[index].Request, _ = request.GetLists(lists[index].RequestIds())
		lists[index].Menu, _ = menu.GetLists(lists[index].MenuIds())
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

func GetRequest(context *admin.Context) (interface{}, error) {
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

	noAssign := make([]request.Model, 0)
	exist := current.RequestIds()
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["name"] = goqu.Op{"like": "%" + param.Keyword + "%"}
		where["action"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	where["id"] = goqu.Op{
		"notin": append(append(exist, admin.OnlyLoginRequestIds...), admin.NoCheckLoginRequestIds...),
	}
	total, err := param.Select("admin_request", &noAssign, where)
	if err != nil {
		return nil, err
	}
	assign := make([]request.Model, 0)
	err = models.Db.From(goqu.T("admin_request").As("req")).
		Select("req.*").
		InnerJoin(
			goqu.T("admin_request_relate").As("rel"),
			goqu.On(goqu.Ex{
				"rel.admin_request_id": goqu.I("req.id"),
			}),
		).
		Where(goqu.Ex{"rel.admin_auth_id": param.Id}).
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

func AssignRequest(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id        int `form:"id" validate:"required"`
		RequestId int `form:"requestId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, requestRelate.Add(param.RequestId, param.Id)
}

func RemoveRequest(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id        int `form:"id" validate:"required"`
		RequestId int `form:"requestId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, requestRelate.Remove(param.RequestId, param.Id)
}

func GetUserGroup(context *admin.Context) (interface{}, error) {
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

	noAssign := make([]userGroup.Model, 0)
	exist := current.UserGroupIds()
	where := make(goqu.Ex)
	if param.Keyword != "" {
		where["name"] = goqu.Op{"like": "%" + param.Keyword + "%"}
	}
	if len(exist) != 0 {
		where["id"] = goqu.Op{"notin": exist}
	}
	total, err := param.Select("admin_user_group", &noAssign, where)
	if err != nil {
		return nil, err
	}
	assign := make([]userGroup.Model, 0)
	err = models.Db.From(goqu.T("admin_user_group").As("ug")).
		Select("ug.*").
		InnerJoin(
			goqu.T("admin_user_group_relate").As("rel"),
			goqu.On(goqu.Ex{
				"rel.admin_user_group_id": goqu.I("ug.id"),
			}),
		).
		Where(goqu.Ex{"rel.admin_auth_id": param.Id}).
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

func AssignUserGroup(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id          int `form:"id" validate:"required"`
		UserGroupId int `form:"userGroupId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, userGroupRelate.Add(param.UserGroupId, param.Id)
}

func RemoveUserGroup(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id          int `form:"id" validate:"required"`
		UserGroupId int `form:"userGroupId" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	return nil, userGroupRelate.Remove(param.UserGroupId, param.Id)
}

func GetMenu(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id int `form:"id" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	current, err := thisModel.GetById(param.Id)
	if err != nil {
		return nil, err
	}

	items := make([]struct {
		menu.Model
		Checked bool `db:"-" json:"checked"`
	}, 0)

	err = models.Db.From("admin_menu").ScanStructs(&items)
	if err != nil {
		return nil, err
	}

	exist := current.MenuIds()

	for index := range items {
		items[index].Checked = inSliceInt(items[index].Id, exist)
	}
	return items, nil
}

func AssignMenu(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Id      int   `form:"id" validate:"required"`
		MenuIds []int `form:"menuIds"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	current, err := thisModel.GetById(param.Id)
	if err != nil {
		return nil, err
	}

	// 清空菜单
	if len(param.MenuIds) == 0 {
		return nil, menuRelate.Remove(0, current.Id)
	}

	exist := current.MenuIds()

	// 删除
	temp := set.IntSliceDifference(exist, param.MenuIds)
	if len(temp) > 0 {
		_ = menuRelate.RemoveMultiple(temp, param.Id)
	}

	// 添加
	temp = set.IntSliceDifference(param.MenuIds, exist)
	if len(temp) > 0 {
		_ = menuRelate.AddMultiple(temp, param.Id)
	}
	return nil, nil
}

func inSliceInt(n int, list []int) bool {
	for _, i := range list {
		if i == n {
			return true
		}
	}
	return false
}
