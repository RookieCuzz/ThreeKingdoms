package controllers

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	model2 "ThreeKingdoms/src/services/login/model"
	"ThreeKingdoms/src/services/login/proto"
	"ThreeKingdoms/src/services/models"
	"ThreeKingdoms/src/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (account *Account) SetupRouter(router *net.RouterStruct) {

	g := router.CreateGroup("account")
	g.AddEventHandler("login", account.Login)

}

func (account *Account) Login(request *net.WsMsgRequestStruct, response *net.WsMsgResponseStruct) {

	//到数据库核对账号密码

	//fmt.Printf("账号密码为%s", request.Body.(proto.LoginClientPacketStruct))
	// 假设 request.Body.MsgContent 是从 JSON 解析出来的
	msgContentJSON, err := json.Marshal(request.Body.MsgContent)
	if err != nil {
		fmt.Println("序列化 MsgContent 失败:", err)
		return
	}

	// 定义目标结构体
	var loginData proto.LoginClientPacketStruct

	// 反序列化为 LoginClientPacketStruct
	err = json.Unmarshal(msgContentJSON, &loginData)
	if err != nil {
		fmt.Println("解析 LoginClientPacketStruct 失败:", err)
		return
	}

	fmt.Println("解析成功:", loginData)

	user := models.User{}
	fmt.Println(gamedatabase.Engine.Tables)
	ok, err := gamedatabase.Engine.Table(&user).Where("username = ?", loginData.Username).Get(&user)
	if err != nil {
		log.Println(err)
	}
	//查无此人
	if !ok {
		fmt.Println("查无此人")
		response.Body.Code = constant.UserNotExist
		return
	}
	//密码错误
	//pwd := utils.Password(loginData.Password, user.Passcode)
	//if pwd != user.Passwd {
	//	fmt.Println("密码错误")
	//	response.Body.Code = constant.PwdIncorrect
	//	return
	//}
	token, err := utils.Award(user.UId)

	fmt.Println("验证成功!放行")
	//密码正确 发个登录成功的包给客户端
	response.Body.Code = constant.OK
	packetStruct := &proto.LoginServerPacketStruct{
		Session:  string(token),
		Username: user.Username,
		Password: user.Passwd,
		UUID:     user.UId,
	}
	response.Body.MsgContent = packetStruct

	//保存用户这次登录的细节
	userLoginHistory := &model2.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginData.Ip,
		Hardware: loginData.Hardware, State: model2.Login,
	}
	gamedatabase.Engine.Table(userLoginHistory).Insert(userLoginHistory)
	//刷新用户最后一次登录的信息
	userLoginLast := &model2.LoginLast{}
	ok, _ = gamedatabase.Engine.Table(&userLoginLast).Where("uid=?", user.UId).Get(userLoginHistory)

	if !ok {
		userLoginLast.LoginTime = time.Now()
		userLoginLast.Ip = loginData.Ip
		userLoginLast.Hardware = loginData.Hardware
		userLoginLast.IsLogout = 0
		userLoginLast.Session = string(token)

		gamedatabase.Engine.Table(userLoginLast).Update(userLoginLast)
	} else {
		userLoginLast.LoginTime = time.Now()
		userLoginLast.Ip = loginData.Ip
		userLoginLast.Hardware = loginData.Hardware
		userLoginLast.IsLogout = 0
		userLoginLast.Session = string(token)
		userLoginLast.UId = user.UId
		gamedatabase.Engine.Table(userLoginLast).Insert(userLoginLast)
	}
	//缓存连接？
	net.WebSocketManager.UserLogin(request.Connection, int64(user.UId), token)
}
