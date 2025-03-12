package game

import (
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/game/controllers"
	"ThreeKingdoms/src/services/game/gameconfig"
)

func Init() {
	gameconfig.Basic.Load()
	initRouter()
}

var Router = net.NewRouter()

func GetRouter() *net.RouterStruct {

	return Router
	//初始化路由 路径

}

func initRouter() {
	controllers.DefaultRoleController.SetupRouter(Router)

}
