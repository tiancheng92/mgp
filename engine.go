package mgp

import "github.com/gin-gonic/gin"

type Engine struct {
	*gin.Engine
	groups           []*RouterGroup
	defaultResponses []*ReturnType
}

func New() *Engine {
	return &Engine{Engine: gin.New()}
}

func (e *Engine) GenerateSwagger() {
	GenerateSwagger("", nil, toGoSwagGroup(e.groups), e.defaultResponses, true)
}

func (e *Engine) Group(relativePath string, handlers ...gin.HandlerFunc) *RouterGroup {
	g := &RouterGroup{routerGroup: e.Engine.Group(relativePath, handlers...), groupName: relativePath}
	e.groups = append(e.groups, g)
	return g
}
