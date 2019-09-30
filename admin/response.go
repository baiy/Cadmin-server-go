package admin

import (
	"encoding/json"
	"fmt"
	"io"
)

// 系统对外响应
type Response struct {
	Status string      `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

func (r Response) Json() []byte {
	ret := make([]byte, 0)
	ret, err := json.Marshal(r)
	if err != nil {
		ret = []byte(
			fmt.Sprintf("{\"status\":\"error\",\"info\":\"响应结果序列化错误(JSON ERR:%s)\",\"data\":\"\"}", err),
		)
	}
	return ret
}

func (r Response) JsonResponse(w io.Writer) error {
	_, err := w.Write(r.Json())
	return err
}

func succeeResponse(data interface{}) *Response {
	return &Response{
		Status: "success",
		Info:   "操作成功",
		Data:   data,
	}
}

func errorResponse(info string, data interface{}) *Response {
	return &Response{
		Status: "error",
		Info:   info,
		Data:   data,
	}
}
