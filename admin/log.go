package admin

import (
	"github.com/baiy/Cadmin-service-go/models/request"
	"github.com/baiy/Cadmin-service-go/models/user"
	"time"
)

type LogContent struct {
	// 用户
	User *user.Model
	// 请求
	Request *request.Model
	// 响应
	Response *Response
	// 响应时间
	Time time.Time
}

var LogCallback func(content LogContent)

// 注册密码生成器
func RegisterLogCallback(callback func(content LogContent)) {
	LogCallback = callback
}
