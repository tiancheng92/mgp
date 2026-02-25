package mgp

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pretty66/websocketproxy"
	"github.com/tiancheng92/mgp/errors"
	"github.com/tiancheng92/mgp/errors/default_error_code"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Context struct {
	*gin.Context
	ws *websocket.Conn
}

// cp 创建一个新的 Context 实例
func newContext(ctx *gin.Context) *Context {
	return &Context{ctx, nil}
}

// do 执行函数，如果上下文没有被中止
func (c *Context) do(f func()) *Context {
	if !c.IsAborted() {
		f()
	}
	return c
}

func (c *Context) Proxy(f func() (*httputil.ReverseProxy, error)) {
	c.do(func() {
		p, err := f()
		if err != nil {
			c.renderValidationError(err)
			return
		}
		p.ServeHTTP(c.Context.Writer, c.Context.Request)
	})
}

// HR 处理函数并渲染响应
func (c *Context) HR(f any) {
	c.do(func() {
		var err error
		var resp any

		// 根据不同的函数类型处理
		switch fn := f.(type) {
		case func():
			fn()
		case func() error:
			err = fn()
		case func() (any, error):
			resp, err = fn()
		default:
			err = errors.WithCode(default_error_code.ErrServer, "invalid function type")
		}
		// 渲染响应
		Response(c.Context, resp, err)
	})
}

// HD 处理函数并渲染响应(文件下载)
func (c *Context) HD(f any, filename string) {
	c.do(func() {
		var err error
		var buf *bytes.Buffer

		// 根据不同的函数类型处理
		switch fn := f.(type) {
		case func() (*bytes.Buffer, error):
			buf, err = fn()
		case func() *bytes.Buffer:
			buf = fn()
		default:
			err = errors.New("invalid function type")
		}
		// 渲染响应
		if err != nil {
			Response(c.Context, nil, err)
		} else {
			ResponseDownloadSteam(c.Context, buf, filename)
		}
	})
}

// HW 处理函数并渲染响应
func (c *Context) HW(f func(*websocket.Conn)) {
	c.do(func() {
		var err error
		c.ws, err = upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.renderValidationError(err)
			return
		}
		defer c.ws.Close()

		f(c.ws)
	})
}

func (c *Context) HWP(f func() (*websocketproxy.WebsocketProxy, error)) {
	c.do(func() {
		wp, err := f()
		if err != nil {
			c.handleWebsocketError(err)
			return
		}
		wp.Proxy(c.Writer, c.Request)
	})
}

func (c *Context) handleWebsocketError(err error) {
	if c.ws == nil {
		var err1 error
		c.ws, err1 = upGrader.Upgrade(c.Writer, c.Request, nil)
		if err1 != nil {
			return
		}
	}

	_ = c.ws.WriteMessage(websocket.TextMessage, StringToBytes(color.New(color.FgHiRed, color.Bold).Sprintln(err.Error())))
	c.ws.Close()
}

// BindBody 绑定请求体到指定结构体
func (c *Context) BindBody(ptr ...any) *Context {
	return c.do(func() {
		for i := range ptr {
			if err := c.ShouldBindBodyWithJSON(ptr[i]); err != nil {
				c.renderValidationError(err)
				return
			}
		}
	})
}

// BindQuery 绑定查询参数到指定结构体
func (c *Context) BindQuery(ptr ...any) *Context {
	return c.do(func() {
		for i := range ptr {
			if err := c.ShouldBindQuery(ptr[i]); err != nil {
				c.renderValidationError(err)
				return
			}
		}
	})
}

// BindParams 绑定 URI 参数到指定结构体
func (c *Context) BindParams(ptr ...any) *Context {
	return c.do(func() {
		for i := range ptr {
			if err := c.ShouldBindUri(ptr[i]); err != nil {
				c.renderValidationError(err)
				return
			}
		}
	})
}

// BindHeader 绑定请求头到指定结构体
func (c *Context) BindHeader(ptr ...any) *Context {
	return c.do(func() {
		for i := range ptr {
			if err := c.ShouldBindHeader(ptr[i]); err != nil {
				c.renderValidationError(err)
				return
			}
		}
	})
}

// BindPaginateQuery 绑定分页查询参数
func (c *Context) BindPaginateQuery(ptr *PaginateQuery) *Context {
	return c.do(func() {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if err != nil || pageSize < 1 {
			pageSize = 20
		}

		ptr.Page = page
		ptr.PageSize = pageSize
		ptr.Order = c.DefaultQuery("order", "")
		ptr.Search = c.DefaultQuery("search", "")
		ptr.Params = c.Request.URL.Query()
	})
}

// renderValidationError 渲染验证错误
func (c *Context) renderValidationError(err error) {
	Response(c.Context, nil, HandleValidationErr(err))
}
