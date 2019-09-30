package admin

import (
	"errors"
	"fmt"
	"strings"
)

type DispatchItem struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Dispatch interface {
	Name() string
	Description() string
	Call(*Context) (interface{}, error)
}

var dispatchers = make(map[string]Dispatch)

func RegisterDispatch(type_ string, dispatcher Dispatch) {
	dispatchers[strings.ToLower(type_)] = dispatcher
}

func GetDispatcher(type_ string) (Dispatch, error) {
	type_ = strings.ToLower(type_)
	dispatcher, is := dispatchers[type_]
	if !is {
		return nil, errors.New(fmt.Sprintf("未找到请求类型(%s)对应的调度程序", type_))
	}
	return dispatcher, nil
}

func AllDispatcher() map[string]Dispatch {
	return dispatchers
}

func AllDispatcherLength() int {
	return len(dispatchers)
}
