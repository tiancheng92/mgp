package mgp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RouterGroup struct {
	routerGroup          *gin.RouterGroup
	groupName            string
	groups               []*RouterGroup
	routes               []*Route
	defaultTags          []string
	defaultAccepts       []string
	defaultProduces      []string
	defaultUseApiKeyAuth bool
}

func (g *RouterGroup) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	sg := &RouterGroup{
		routerGroup:          g.routerGroup.Group(relativePath, handlers...),
		groupName:            getFullPath(g.groupName, relativePath),
		defaultTags:          g.defaultTags,
		defaultAccepts:       g.defaultAccepts,
		defaultProduces:      g.defaultProduces,
		defaultUseApiKeyAuth: g.defaultUseApiKeyAuth,
	}
	g.groups = append(g.groups, sg)
	return sg
}

func (g *RouterGroup) SetTagsForSwagger(tags ...string) *RouterGroup {
	g.defaultTags = tags
	return g
}

func (g *RouterGroup) SetAcceptsForSwagger(accepts ...string) *RouterGroup {
	g.defaultAccepts = accepts
	return g
}

func (g *RouterGroup) SetProducesForSwagger(produces ...string) *RouterGroup {
	g.defaultProduces = produces
	return g
}

func (g *RouterGroup) SetUseApiKeyAuthForSwagger() *RouterGroup {
	g.defaultUseApiKeyAuth = true
	return g
}

func (g *RouterGroup) Handle(httpMethod, relativePath string, f func(c *Context)) Swagger {
	g.routerGroup.Handle(httpMethod, relativePath, func(c *gin.Context) {
		f(newContext(c))
	})
	gl := strings.Split(g.groupName, "/")
	gr := &Route{
		Path:          getFullPath(g.groupName, relativePath),
		Method:        httpMethod,
		FuncName:      fmt.Sprintf("%s%s", camelString(gl[len(gl)-1]), getFuncName(f)),
		Tags:          g.defaultTags,
		Accepts:       g.defaultAccepts,
		Produces:      g.defaultProduces,
		UseApiKeyAuth: g.defaultUseApiKeyAuth,
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
