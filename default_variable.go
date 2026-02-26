package mgp

var (
	successMsg   = "Success"
	successCode  = 200000
	swagFilePath = "./goswag.go"
)

func SetSuccessMsg(msg string) {
	successMsg = msg
}

func SetSuccessCode(code int) {
	successCode = code
}

func SetSwagFileName(path string) {
	swagFilePath = path
}
