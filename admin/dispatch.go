package admin

import (
	"errors"
	"fmt"
	"strings"
)

type Dispatch interface {
	Name() string
	Description() string
	Call(*Context) (interface{}, error)
}

var dispatchers = make(map[string]Dispatch)

func RegisterDispatch(name string, dispatcher Dispatch) {
	dispatchers[strings.ToLower(name)] = dispatcher
}

func GetDispatcher(name string) (Dispatch, error) {
	dispatcher, is := dispatchers[name]
	if !is {
		return nil, errors.New(fmt.Sprintf("未找到请求类型(%s)对应的调度程序", name))
	}
	return dispatcher, nil
}

func AllDispatcher() map[string]Dispatch {
	return dispatchers
}
