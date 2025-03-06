package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/gamedatabase"
	"ThreeKingdoms/src/services/web"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	//配置初始化
	config.Init()
	//数据库初始化
	gamedatabase.TestDatabase()
	router := gin.Default()
	web.Init(router)
	serverConfig := config.Config.WebServer
	router.Run(serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port))
}
