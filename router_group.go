package mgp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	routerGroup *gin.RouterGroup
	groupName   string
	groups      []*RouterGroup
	routes      []*Route
}

func (g *RouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	sg := &RouterGroup{routerGroup: g.routerGroup.Group(relativePath, handlers...), groupName: getFullPath(g.groupName, relativePath)}
	g.groups = append(g.groups, sg)
	return sg
}

func (g *RouterGroup) Handle(httpMethod, relativePath string, f func(c *Context)) Swagger {
	g.routerGroup.Handle(httpMethod, relativePath, func(c *gin.Context) {
		f(newContext(c))
	})
	gl := strings.Split(g.groupName, "/")
	gr := &Route{
		Path:     getFullPath(g.groupName, relativePath),
		Method:   httpMethod,
		FuncName: fmt.Sprintf("%s%s", camelString(gl[len(gl)-1]), getFuncName(f)),
	}
	g.routes = append(g.routes, gr)
	return gr
}

func (g *RouterGroup) GET(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodGet, relativePath, f)
}

func (g *RouterGroup) POST(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodPost, relativePath, f)
}

func (g *RouterGroup) PUT(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodPut, relativePath, f)
}

func (g *RouterGroup) PATCH(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodPatch, relativePath, f)
}

func (g *RouterGroup) DELETE(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodDelete, relativePath, f)
}

func (g *RouterGroup) Options(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodOptions, relativePath, f)
}

func (g *RouterGroup) Head(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodHead, relativePath, f)
}

func (g *RouterGroup) Connect(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodConnect, relativePath, f)
}

func (g *RouterGroup) Trace(relativePath string, f func(c *Context)) Swagger {
	return g.Handle(http.MethodTrace, relativePath, f)
}
