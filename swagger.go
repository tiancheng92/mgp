package mgp

type Swagger interface {
	SetSummaryForSwagger(summary string) Swagger
	SetDescriptionForSwagger(description string) Swagger
	SetTagsForSwagger(tags ...string) Swagger
	SetAcceptsForSwagger(accept ...string) Swagger
	SetProducesForSwagger(produce ...string) Swagger
	SetBodyForSwagger(data interface{}) Swagger
	SetQueryForSwagger(query interface{}) Swagger
	SetPathForSwagger(path interface{}) Swagger
	SetHeaderForSwagger(header interface{}) Swagger
	SetReturnsForSwagger(data ...*ReturnType) Swagger
	SetUseApiKeyAuthForSwagger() Swagger
}
