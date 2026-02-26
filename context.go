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

func newContext(ctx *gin.Context) *Context {
	return &Context{ctx, nil}
}

func (c *Context) do(f func()) *Context {
	if !c.IsAborted() {
		f()
	}
	return c
}

func (c *Context) HP(f func() (*httputil.ReverseProxy, error)) {
	c.do(func() {
		p, err := f()
		if err != nil {
			c.renderValidationError(err)
			return
		}
		p.ServeHTTP(c.Context.Writer, c.Context.Request)
	})
}

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

func (c *Context) renderValidationError(err error) {
	Response(c.Context, nil, HandleValidationErr(err))
}
