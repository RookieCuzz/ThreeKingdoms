package gateway

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/gateway/controllers"
)

func Init() {

	initRouter()
}

var Router = net.NewRouter()

func GetRouter() *net.RouterStruct {

	return Router
	//初始化路由 路径

}

func initRouter() {
	group := Router.CreateGroup("*")
	group.AddEventHandler("*", controllers.Gateway.All)
	proxyConfig := config.Config.GateServer

	controllers.Gateway.LoginProxy = proxyConfig.LoginProxy
	controllers.Gateway.GameProxy = proxyConfig.GameProxy

}
