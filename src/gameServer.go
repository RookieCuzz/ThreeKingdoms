package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/net"
	"ThreeKingdoms/src/services/game"
	"strconv"
	"sync"
)

/**


1.登录完成了,创建角色（玩家 ）
2.需要根据用户,查询其拥有的角色,没有则创建角色
3.木材,铁 令牌 金钱初始化
4.地图 城池 要塞需要定义
5.资源,军队,城池,武将
*/

func main() {
	var wg sync.WaitGroup
	//配置初始化
	config.Init()
	//数据库初始化
	gamedatabase.TestDatabase()
	port := config.Config.GameServer.Port
	host := config.Config.GameServer.Host
	wg.Add(1)

	//构建服务器的路由树
	game.Init()

	//开启服务器
	server := net.NewServer(host + ":" + strconv.Itoa(port))
	//设置服务器的路由树
	server.Router = game.GetRouter()
	//启动登录服务器
	server.Start()
}
