package proto

//登录客户端包
type LoginClientPacketStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
	UUID     string `json:"uuid"`
}

//登录服务端包
type LoginServerPacketStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Hardware string `json:"hardware"`
}
