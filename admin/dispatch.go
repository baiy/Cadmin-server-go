package admin

import (
	"errors"
	"fmt"
	"strings"
)

// 调度器接口
type Dispatch interface {
	// 调度器标识
	Key() string
	// 调度器名称
	Name() string
	// 调度器描述
	Description() string
	// 请求调度方法
	Call(*Context) (interface{}, error)
}

var dispatchers = make(map[string]Dispatch)

func RegisterDispatch(dispatcher Dispatch) {
	dispatchers[strings.ToLower(dispatcher.Key())] = dispatcher
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

// 默认调度器
type defaultDispatcher struct {
	HandleMethod map[string]DefaultDispatcherHandleMethod
}

func (d *defaultDispatcher) Key() string {
	return "default"
}

func (d *defaultDispatcher) Name() string {
	return "默认"
}

func (d *defaultDispatcher) Description() string {
	return "Cadmin系统内置的默认请求调度器"
}

func (d *defaultDispatcher) Call(c *Context) (interface{}, error) {
	method, is := d.HandleMethod[c.Request.Call]
	if !is {
		return nil, errors.New(fmt.Sprintf("[%s 未注册对应的处理方法]", c.Request.Action))
	}
	r, err := method(c)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (d *defaultDispatcher) Register(methods map[string]DefaultDispatcherHandleMethod) {
	for name, method := range methods {
		if _, is := d.HandleMethod[name]; is {
			panic("[%s] 对应的处理方法已经存在")
		}
		d.HandleMethod[name] = method
	}
}

type DefaultDispatcherHandleMethod func(*Context) (interface{}, error)

var DefaultDispatcher = &defaultDispatcher{HandleMethod: make(map[string]DefaultDispatcherHandleMethod)}

// 注册默认调度器请求处理方法
func RegisterDefaultDispatcherHandleMethod(methods map[string]DefaultDispatcherHandleMethod) {
	DefaultDispatcher.Register(methods)
}

func init() {
	RegisterDispatch(DefaultDispatcher)
}
