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

func (r *Route) SetSummaryForSwagger(summary string) Swagger {
	r.Summary = summary
	return r
}

func (r *Route) SetDescriptionForSwagger(description string) Swagger {
	r.Description = description
	return r
}

func (r *Route) SetTagsForSwagger(tags ...string) Swagger {
	r.Tags = tags
	return r
}

func (r *Route) SetAcceptsForSwagger(accepts ...string) Swagger {
	r.Accepts = accepts
	return r
}

func (r *Route) SetProducesForSwagger(produces ...string) Swagger {
	r.Produces = produces
	return r
}

func (r *Route) SetBodyForSwagger(body interface{}) Swagger {
	r.BodyStruct = body
	return r
}

func (r *Route) SetQueryForSwagger(query interface{}) Swagger {
	r.QueryStruct = query
	return r
}

func (r *Route) SetPathForSwagger(path interface{}) Swagger {
	r.PathStruct = path
	return r
}

func (r *Route) SetHeaderForSwagger(header interface{}) Swagger {
	r.HeaderStruct = header
	return r
}

func (r *Route) SetReturnsForSwagger(returns ...*ReturnType) Swagger {
	r.Returns = returns
	return r
}

func (r *Route) SetUseApiKeyAuthForSwagger() Swagger {
	r.UseApiKeyAuth = true
	return r
}
