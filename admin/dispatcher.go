// 系统默认调度器
package admin

import (
	"errors"
	"fmt"
)

type dispatcher struct {
	HandleMethod map[string]DefaultDispatcherHandleMethod
}

type DefaultDispatcherHandleMethod func(*Context) (interface{}, error)

var DefaultDispatcher = &dispatcher{HandleMethod: make(map[string]DefaultDispatcherHandleMethod)}

func (d *dispatcher) Name() string {
	return "默认"
}

func (d *dispatcher) Description() string {
	return "Cadmin系统内置的默认请求调度器"
}

func (d *dispatcher) Call(c *Context) (interface{}, error) {
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

func (d *dispatcher) Register(methods map[string]DefaultDispatcherHandleMethod) {
	for name, method := range methods {
		if _, is := d.HandleMethod[name]; is {
			panic("[%s] 对应的处理方法已经存在")
		}
		d.HandleMethod[name] = method
	}
}

// 注册默认调度器请求处理方法
func RegisterDefaultDispatcherHandleMethod(methods map[string]DefaultDispatcherHandleMethod) {
	DefaultDispatcher.Register(methods)
}

func init() {
	RegisterDispatch("default", DefaultDispatcher)
}
