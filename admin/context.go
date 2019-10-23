package admin

import (
	"errors"
	"fmt"
	"github.com/baiy/Cadmin-server-go/models/request"
	"github.com/baiy/Cadmin-server-go/models/token"
	"github.com/baiy/Cadmin-server-go/models/user"
	"github.com/go-playground/form"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zhTranslations "gopkg.in/go-playground/validator.v9/translations/zh"
	"net/http"
	"strconv"
	"time"
)

var (
	validate *validator.Validate
	decoder  *form.Decoder
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

// 请求标识变量名
var ActionName = "_action"

// 请求token变量名
var TokenName = "_token"

// 无需登录的请求ID
var NoCheckLoginRequestIds = []int{1}

// 无需检查权限/只需要登录的请求ID
var OnlyLoginRequestIds = []int{2, 3, 4}

// 添加无需登录的请求ID
func AddNoCheckLoginRequestId(ids ...int) {
	NoCheckLoginRequestIds = append(NoCheckLoginRequestIds, ids...)
}

// 添加无需检查权限/只需要登录的请求ID
func AddOnlyLoginRequestId(ids ...int) {
	OnlyLoginRequestIds = append(OnlyLoginRequestIds, ids...)
}

func init() {
	zhTranslator := zh.New()
	uni = ut.New(zhTranslator, zhTranslator)
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
	decoder = form.NewDecoder()
}

type Context struct {
	Action             string
	Token              string
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
	Request            *request.Model
	User               *user.Model
	Response           *Response
}

func (c *Context) run() {
	defer func() {
		if r := recover(); r != nil {
			c.SetResponse(errorResponse(fmt.Sprint(r), nil))
		}
	}()

	if err := c.initRequest(); err != nil {
		c.SetResponse(errorResponse(err.Error(), nil))
		return
	}

	_ = c.initUser()

	if err := c.checkAccess(); err != nil {
		c.SetResponse(errorResponse(err.Error(), nil))
		return
	}
	result, err := c.call()
	if err != nil {
		c.SetResponse(errorResponse(err.Error(), result))
		return
	}
	c.SetResponse(succeeResponse(result))
	return
}

func (c *Context) SetResponse(r *Response) {
	if LogCallback != nil {
		LogCallback(LogContent{User: c.User, Request: c.Request, Response: c.Response, Time: time.Now()})
	}
	c.Response = r
}

func (c *Context) Output() error {
	return c.Response.JsonResponse(c.HttpResponseWriter)
}

func (c *Context) initRequest() error {
	c.Action = c.HttpRequest.URL.Query().Get(ActionName)
	if c.Action == "" {
		return errors.New("action参数错误")
	}
	req, err := request.GetByAction(c.Action)
	if err != nil {
		return err
	}
	c.Request = req
	return nil
}

func (c *Context) initUser() error {
	c.Token = c.HttpRequest.URL.Query().Get(TokenName)
	if c.Token == "" {
		return errors.New("token 为空")
	}
	req, err := token.GetByToken(c.Token)
	if err != nil {
		return err
	}
	if req.IsExpire() {
		return errors.New("登录状态已过期")
	}

	c.User, err = user.GetById(req.AdminUserId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) checkAccess() error {
	if inSliceInt(c.Request.Id, NoCheckLoginRequestIds) {
		return nil
	}
	if c.User == nil {
		return errors.New("未登录系统")
	}

	if c.User.IsDisabled() {
		return errors.New("用户已被禁用")
	}
	if inSliceInt(c.Request.Id, OnlyLoginRequestIds) {
		return nil
	}

	err := user.CheckAuth(c.User.Id, c.Request.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) call() (interface{}, error) {
	dispatcher, err := GetDispatcher(c.Request.Type)
	if err != nil {
		return nil, err
	}
	return dispatcher.Call(c)
}

func (c *Context) InputInt(name string, def ...int) (int, error) {
	strv := c.HttpRequest.Form.Get(name)
	if len(strv) == 0 && len(def) > 0 {
		return def[0], nil
	}
	return strconv.Atoi(strv)
}

func (c *Context) Input(name string, def ...string) string {
	if v := c.HttpRequest.Form.Get(name); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (c *Context) Form(parameter interface{}) error {
	err := decoder.Decode(parameter, c.HttpRequest.Form)
	if err != nil {
		return err
	}
	err = validate.Struct(parameter)
	if err != nil {
		for _, value := range err.(validator.ValidationErrors).Translate(trans) {
			return errors.New(value)
		}
	}
	return nil
}

// 请求入口方法
func NewContext(rw http.ResponseWriter, r *http.Request) *Context {
	_ = r.ParseForm()
	_ = r.ParseMultipartForm(128)
	context := &Context{HttpResponseWriter: rw, HttpRequest: r}
	context.run()
	return context
}

func inSliceInt(n int, list []int) bool {
	for _, i := range list {
		if i == n {
			return true
		}
	}
	return false
}
