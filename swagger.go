package mgp

type Swagger interface {
	SwaggerSummary(summary string) Swagger
	SwaggerDescription(description string) Swagger
	SwaggerTags(tags ...string) Swagger
	SwaggerAccepts(accept ...string) Swagger
	SwaggerProduces(produce ...string) Swagger
	SwaggerBody(data interface{}) Swagger
	SwaggerQuery(query interface{}) Swagger
	SwaggerPath(path interface{}) Swagger
	SwaggerHeader(header interface{}) Swagger
	SwaggerReturns(data ...*ReturnType) Swagger
	SwaggerUseApiKeyAuth() Swagger
	SwaggerHidde() Swagger
}
