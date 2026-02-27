package mgp

type Route struct {
	Path          string
	Method        string
	FuncName      string
	Summary       string
	Description   string
	Tags          []string
	Accepts       []string
	Produces      []string
	BodyStruct    any
	QueryStruct   any
	HeaderStruct  any
	PathStruct    any
	Returns       []*ReturnType
	UseApiKeyAuth bool
	Hidde         bool
}

func (r *Route) SwaggerSummary(summary string) Swagger {
	r.Summary = summary
	return r
}

func (r *Route) SwaggerDescription(description string) Swagger {
	r.Description = description
	return r
}

func (r *Route) SwaggerTags(tags ...string) Swagger {
	r.Tags = tags
	return r
}

func (r *Route) SwaggerAccepts(accepts ...string) Swagger {
	r.Accepts = accepts
	return r
}

func (r *Route) SwaggerProduces(produces ...string) Swagger {
	r.Produces = produces
	return r
}

func (r *Route) SwaggerBody(body interface{}) Swagger {
	r.BodyStruct = body
	return r
}

func (r *Route) SwaggerQuery(query interface{}) Swagger {
	r.QueryStruct = query
	return r
}

func (r *Route) SwaggerPath(path interface{}) Swagger {
	r.PathStruct = path
	return r
}

func (r *Route) SwaggerHeader(header interface{}) Swagger {
	r.HeaderStruct = header
	return r
}

func (r *Route) SwaggerReturns(returns ...*ReturnType) Swagger {
	r.Returns = returns
	return r
}

func (r *Route) SwaggerUseApiKeyAuth() Swagger {
	r.UseApiKeyAuth = true
	return r
}

func (r *Route) SwaggerHidde() Swagger {
	r.Hidde = true
	return r
}
