package net

type RequestStruct struct {

	//Seq用来匹配 请求和响应
	Seq int64 `json:"seq"`
	//用来区别消息类型 登录消息,充值消息etc
	Name string `json:"name"`
	//消息内容 可以是任何类型
	MsgContent interface{} `json:"msg"`
	//代理服务器
	Proxy string `json:"proxy"`
}

type ResponseStruct struct {
	//Seq用来匹配 请求和响应
	Seq int64 `json:"seq"`
	//用来区别消息类型 登录消息,充值消息etc
	Name string `json:"name"`
	//消息内容 可以是任何类型
	MsgContent interface{} `json:"msg"`
	//响应码 模仿 http
	Code int `json:"code"`
}

type WsMsgRequestStruct struct {
	Body       *RequestStruct `json:"body"`
	Connection WsConnection
}

type WsMsgResponseStruct struct {
	Body *ResponseStruct `json:"body"`
}

// request请求会有参数,往请求里面放参数
type WsConnection interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Heartbeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}
