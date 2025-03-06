package common

type ResultStruct struct {
	Code   int         `json:"code"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func Error(code int, msg string) *ResultStruct {
	return &ResultStruct{
		Code:   code,
		Errmsg: msg,
	}
}

func Success(code int, data interface{}) *ResultStruct {
	return &ResultStruct{
		Code: code,
		Data: data,
	}
}
