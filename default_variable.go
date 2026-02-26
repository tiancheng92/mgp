package mgp

var (
	SuccessMsg   = "Success"
	SuccessCode  = 200000
	SwagFilePath = "./goswag.go"
)

func SetSuccessMsg(msg string) {
	SuccessMsg = msg
}

func SetSuccessCode(code int) {
	SuccessCode = code
}

func SetSwagFileName(path string) {
	SwagFilePath = path
}
