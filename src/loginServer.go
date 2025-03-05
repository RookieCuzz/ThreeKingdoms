package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services"
	"ThreeKingdoms/src/services/controllers"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	//配置初始化
	config.Init()
	//数据库初始化
	gamedatabase.TestDatabase()
	port := config.Config.LoginServer.Port
	host := config.Config.LoginServer.Host
	wg.Add(1)

	//构建服务器的路由树
	router := services.GetRouter()
	//将账号登录添加进路由
	controllers.DefaultAccount.SetupRouter(router)

	//开启服务器
	server := net.NewServer(host + ":" + strconv.Itoa(port))
	//设置服务器的路由树
	server.Router = router
	//启动登录服务器
	server.Start()

}
