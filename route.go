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
}

func (r *Route) SetSummary(summary string) Swagger {
	r.Summary = summary
	return r
}

func (r *Route) SetDescription(description string) Swagger {
	r.Description = description
	return r
}

func (r *Route) SetTags(tags ...string) Swagger {
	r.Tags = tags
	return r
}

func (r *Route) SetAccepts(accepts ...string) Swagger {
	r.Accepts = accepts
	return r
}

func (r *Route) SetProduces(produces ...string) Swagger {
	r.Produces = produces
	return r
}

func (r *Route) SetBody(body interface{}) Swagger {
	r.BodyStruct = body
	return r
}

func (r *Route) SetQuery(query interface{}) Swagger {
	r.QueryStruct = query
	return r
}

func (r *Route) SetPath(path interface{}) Swagger {
	r.PathStruct = path
	return r
}

func (r *Route) SetHeader(header interface{}) Swagger {
	r.HeaderStruct = header
	return r
}

func (r *Route) SetReturns(returns ...*ReturnType) Swagger {
	r.Returns = returns
	return r
}

func (r *Route) SetUseApiKeyAuth() Swagger {
	r.UseApiKeyAuth = true
	return r
}
