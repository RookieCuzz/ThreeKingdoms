package controllers

import (
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/proto"
	"fmt"
	"strconv"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (account *Account) SetupRouter(router *net.RouterStruct) {

	g := router.CreateGroup("account")
	g.AddEventHandler("login", account.Login)

}

func (account *Account) Login(request *net.WsMsgRequestStruct, response *net.WsMsgResponseStruct) {
	content := request.Body.MsgContent
	fmt.Println("账号密码为" + content.(string))

	packetStruct := &proto.LoginClientPacketStruct{
		Username: "聪明的小红",
		UUID:     strconv.Itoa(11455),
		Password: "123456xxx",
		Session:  "testsession",
	}
	response.Body.MsgContent = packetStruct
	//response.MsgContent=&proto.LoginServerPacketStruct{
	//	Username: "聪明的小红",
	//	Password: "123456xxx",
	//	Ip:       "127.0.0.1",
	//	Hardware: "Windows",
	//
	//}
}
