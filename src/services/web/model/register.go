package model

type RegisterRequestStruct struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Hardware string `form:"hardware" json:"hardware"`
}
