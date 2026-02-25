package mgp

import (
	"bytes"
	"fmt"
	"mgp/errors"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

type Result[D any] struct {
	Data D      `json:"data"` // 返回数据/错误详细信息
	Msg  string `json:"msg"`  // 请求结果
	Code int    `json:"code"` // 状态码
}

func Response(ctx *gin.Context, data any, err error) {
	if ctx.IsAborted() {
		return
	}
	ctx.Abort()

	result := new(Result[any])
	if err != nil {
		handleError(ctx, result, err)
	} else {
		handleSuccess(ctx, result, data)
	}
}

func handleError(ctx *gin.Context, result *Result[any], err error) {
	coder := errors.ParseCoder(err)
	result.Code = coder.Code()
	if coder.String() != "" {
		result.Msg = coder.String()
	} else {
		result.Msg = "Unknown Error Type"
	}

	ctx.Set(AllError, err)
	if coder.HTTPStatus() >= 400 && coder.HTTPStatus() < 500 {
		ctx.Set(ErrorLogLevelWarn, err)
	} else {
		ctx.Set(ErrorLogLevelError, err)
	}
	result.Data = err.Error()
	b, _ := sonic.Marshal(result)
	ctx.Data(coder.HTTPStatus(), "application/json", b)
}

func handleSuccess(ctx *gin.Context, result *Result[any], data any) {
	if d, ok := data.(PaginateInterface); ok {
		pd := new(PaginateData[any])
		pd.Items = d.GetItems()
		pd.Paginate = d.GetPaginate()
		result.Data = pd
	} else {
		result.Data = data
	}
	result.Msg = "Success"
	result.Code = 200000
	b, _ := sonic.Marshal(result)
	ctx.Data(http.StatusOK, "application/json", b)
}

func ResponseDownloadSteam(ctx *gin.Context, buf *bytes.Buffer, fileName string) {
	if ctx.IsAborted() {
		return
	}
	ctx.Abort()

	ctx.Header("Content-Length", strconv.Itoa(buf.Len()))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	p := make([]byte, 4096)

	for {
		count, err := buf.Read(p)
		if err != nil {
			break
		}
		if count > 0 {
			if _, err := ctx.Writer.Write(p[:count]); err != nil {
				break
			}
			ctx.Writer.Flush()
		}
	}
}
