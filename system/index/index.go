package index

import (
	"errors"
	"github.com/baiy/Cadmin-server-go/admin"
	"github.com/baiy/Cadmin-server-go/models/auth"
	"github.com/baiy/Cadmin-server-go/models/menu"
	"github.com/baiy/Cadmin-server-go/models/request"
	"github.com/baiy/Cadmin-server-go/models/token"
	"github.com/baiy/Cadmin-server-go/models/user"
	"github.com/baiy/Cadmin-server-go/models/userGroup"
	"net"
	"net/http"
	"strings"
)

func Login(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Username string `form:"username" validate:"required"`
		Password string `form:"password" validate:"required"`
	})
	err := context.Form(param)
	if err != nil {
		return nil, err
	}

	u, err := user.GetByUserName(param.Username)
	if err != nil {
		return nil, err
	}

	if u.IsDisabled() {
		return nil, errors.New("用户已经禁用")
	}

	if !admin.Passworder.Verify([]byte(param.Password), []byte(u.Password)) {
		return nil, errors.New("密码错误")
	}

	// 清理token
	token.Clear()

	// 添加token
	t := token.Add(u.Id)
	// 更新用户登陆
	u.LoginUpdate(clientIP(context.HttpRequest))

	return map[string]string{"token": t}, nil
}

func Logout(context *admin.Context) (interface{}, error) {
	if context.Token == "" {
		return nil, errors.New("token 错误")
	}
	token.Remove(context.Token)
	return nil, nil
}

func Load(context *admin.Context) (interface{}, error) {
	allUser, err := user.All()
	if err != nil {
		return nil, err
	}
	menus, err := menu.GetLists(context.User.MenuIds())
	if err != nil {
		return nil, err
	}
	requests, err := request.GetLists(context.User.RequestIds())
	if err != nil {
		return nil, err
	}
	auths, err := auth.GetLists(context.User.AuthIds())
	if err != nil {
		return nil, err
	}
	userGroups, err := userGroup.GetLists(context.User.UserGroupIds())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"user":      context.User,
		"allUser":   allUser,
		"menu":      menus,
		"request":   requests,
		"auth":      auths,
		"userGroup": userGroups,
	}, nil
}

func CurrentSetting(context *admin.Context) (interface{}, error) {
	param := new(struct {
		Username       string `form:"username" validate:"required"`
		Password       string `form:"password"`
		RepeatPassword string `form:"repeatPassword"`
	})

	err := context.Form(param)
	if err != nil {
		return nil, err
	}
	password := ""
	if param.Password != "" {
		if param.Password != param.RepeatPassword {
			return nil, errors.New("两次输入密码不一致")
		}
		password = string(admin.Passworder.Hash([]byte(param.Password)))
	}
	return nil, user.Updata(context.User.Id, param.Username, password, 0)
}

func clientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
