// 请求处理方法注册

package system

import (
	"github.com/baiy/Cadmin-service-go/admin"
	"github.com/baiy/Cadmin-service-go/system/auth"
	"github.com/baiy/Cadmin-service-go/system/index"
	"github.com/baiy/Cadmin-service-go/system/menu"
	"github.com/baiy/Cadmin-service-go/system/request"
	"github.com/baiy/Cadmin-service-go/system/user"
	"github.com/baiy/Cadmin-service-go/system/userGroup"
)

func init() {
	admin.RegisterDefaultDispatcherHandleMethod(map[string]admin.DefaultDispatcherHandleMethod{
		"Baiy.Cadmin.System.Index.login":          index.Login,
		"Baiy.Cadmin.System.Index.logout":         index.Logout,
		"Baiy.Cadmin.System.Index.load":           index.Load,
		"Baiy.Cadmin.System.User.lists":           user.Lists,
		"Baiy.Cadmin.System.User.save":            user.Save,
		"Baiy.Cadmin.System.User.remove":          user.Remove,
		"Baiy.Cadmin.System.UserGroup.lists":      userGroup.Lists,
		"Baiy.Cadmin.System.UserGroup.save":       userGroup.Save,
		"Baiy.Cadmin.System.UserGroup.remove":     userGroup.Remove,
		"Baiy.Cadmin.System.UserGroup.getUser":    userGroup.GetUser,
		"Baiy.Cadmin.System.UserGroup.assignUser": userGroup.AssignUser,
		"Baiy.Cadmin.System.UserGroup.removeUser": userGroup.RemoveUser,
		"Baiy.Cadmin.System.Request.lists":        request.Lists,
		"Baiy.Cadmin.System.Request.save":         request.Save,
		"Baiy.Cadmin.System.Request.remove":       request.Remove,
		"Baiy.Cadmin.System.Request.type":         request.Type,
		"Baiy.Cadmin.System.Menu.lists":           menu.Lists,
		"Baiy.Cadmin.System.Menu.sort":            menu.Sort,
		"Baiy.Cadmin.System.Menu.save":            menu.Save,
		"Baiy.Cadmin.System.Menu.remove":          menu.Remove,
		"Baiy.Cadmin.System.Auth.lists":           auth.Lists,
		"Baiy.Cadmin.System.Auth.save":            auth.Save,
		"Baiy.Cadmin.System.Auth.remove":          auth.Remove,
		"Baiy.Cadmin.System.Auth.getRequest":      auth.GetRequest,
		"Baiy.Cadmin.System.Auth.assignRequest":   auth.AssignRequest,
		"Baiy.Cadmin.System.Auth.removeRequest":   auth.RemoveRequest,
		"Baiy.Cadmin.System.Auth.getUserGroup":    auth.GetUserGroup,
		"Baiy.Cadmin.System.Auth.assignUserGroup": auth.AssignUserGroup,
		"Baiy.Cadmin.System.Auth.removeUserGroup": auth.RemoveUserGroup,
		"Baiy.Cadmin.System.Auth.getMenu":         auth.GetMenu,
		"Baiy.Cadmin.System.Auth.assignMenu":      auth.AssignMenu,
	})
}
