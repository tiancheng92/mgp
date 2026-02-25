package mgp

type ErrorLogLevel string

const (
	ErrorLogLevelWarn  ErrorLogLevel = "warn"
	ErrorLogLevelError ErrorLogLevel = "error"
)

const AllError = "all_error"
