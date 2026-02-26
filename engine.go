package mgp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	*gin.Engine
	routes           []*Route
	groups           []*RouterGroup
	defaultResponses []*ReturnType
}

func New(defaultResponses ...*ReturnType) *Engine {
	return &Engine{Engine: gin.New(), defaultResponses: defaultResponses}
}

func (e *Engine) GenerateSwagger() {
	GenerateSwagger(e.routes, toGoSwagGroup(e.groups), e.defaultResponses, true)
}

func (e *Engine) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	g := &RouterGroup{routerGroup: e.Engine.Group(relativePath, handlers...), groupName: relativePath}
	e.groups = append(e.groups, g)
	return g
}

func (e *Engine) Handle(httpMethod, relativePath string, f func(c *Context)) Swagger {
	e.Engine.Handle(httpMethod, relativePath, func(c *gin.Context) {
		f(newContext(c))
	})
	gr := &Route{
		Path:     relativePath,
		Method:   httpMethod,
		FuncName: getFuncName(f),
	}
	e.routes = append(e.routes, gr)
	return gr
}

func (e *Engine) GET(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodGet, relativePath, f)
}

func (e *Engine) POST(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodPost, relativePath, f)
}

func (e *Engine) PUT(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodPut, relativePath, f)
}

func (e *Engine) PATCH(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodPatch, relativePath, f)
}

func (e *Engine) DELETE(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodDelete, relativePath, f)
}

func (e *Engine) Options(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodOptions, relativePath, f)
}

func (e *Engine) Head(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodHead, relativePath, f)
}

func (e *Engine) Connect(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodConnect, relativePath, f)
}

func (e *Engine) Trace(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodTrace, relativePath, f)
}
