package mgp

type Swagger interface {
	SetSummary(summary string) Swagger
	SetDescription(description string) Swagger
	SetTags(tags ...string) Swagger
	SetAccepts(accept ...string) Swagger
	SetProduces(produce ...string) Swagger
	SetBody(data interface{}) Swagger
	SetQuery(query interface{}) Swagger
	SetPath(path interface{}) Swagger
	SetHeader(header interface{}) Swagger
	SetReturns(data ...*ReturnType) Swagger
	SetUseApiKeyAuth() Swagger
}
