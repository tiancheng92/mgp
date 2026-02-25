package default_error_code

import "github.com/tiancheng92/mgp/errors"

const ErrClientParam string = "ErrClientParam"
const ErrServer string = "ErrServer"

func init() {
	errors.Register(ErrClientParam, 400000, 400, "参数异常")
	errors.Register(ErrServer, 500000, 500, "服务端错误")
}
