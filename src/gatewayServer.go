package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/gateway"
	"fmt"
	"strconv"
)

/*
	1.登录功能 account.login 需要通过网关 转发登录服务器
	2. 网关再作为客户端 与登录服务器的websocket进行交互
	3. 网关又和游戏客户端进行交互,网关是websocket的服务端
	4.websocket的服务端已经实现
	5.路由 gateway服务端接受所有进行请求(*)
	8.握手协议 检测第一次建立连接的时候进行授权访问

	//也就是说gateway同时要维护 websocket服务端和websocket客户端


*/

func main() {

	//配置初始化
	config.Init()
	fmt.Printf("99999999")
	//数据库初始化
	gamedatabase.TestDatabase()
	port := config.Config.GateServer.Port
	host := config.Config.GateServer.Host
	//初始化路由和对应处理器
	gateway.Init()

	//开启服务器
	server := net.NewServer(host + ":" + strconv.Itoa(port))
	fmt.Println(host + ":" + strconv.Itoa(port))
	//设置服务器的路由树
	server.Router = gateway.GetRouter()
	//启动登录服务器
	server.Start()

}
