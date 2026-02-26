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

func (e *Engine) Use(handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.Use(handlers...)
}

func (e *Engine) NoRoute(handlers ...gin.HandlerFunc) {
	e.Engine.NoRoute(handlers...)
}

func (e *Engine) NoMethod(handlers ...gin.HandlerFunc) {
	e.Engine.NoMethod(handlers...)
}

func (e *Engine) RawGET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.GET(relativePath, handlers...)
}

func (e *Engine) RawPOST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.POST(relativePath, handlers...)
}

func (e *Engine) RawPUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.PUT(relativePath, handlers...)
}

func (e *Engine) RawPATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.PATCH(relativePath, handlers...)
}

func (e *Engine) RawDELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.DELETE(relativePath, handlers...)
}

func (e *Engine) RawOPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.OPTIONS(relativePath, handlers...)
}

func (e *Engine) RawHEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return e.Engine.HEAD(relativePath, handlers...)
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

func (e *Engine) OPTIONS(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodOptions, relativePath, f)
}

func (e *Engine) HEAD(relativePath string, f func(c *Context)) Swagger {
	return e.Handle(http.MethodHead, relativePath, f)
}
