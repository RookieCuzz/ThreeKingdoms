package controllers

import (
	"ThreeKingdoms/src/net"
	"log"
	"strings"
	"sync"
)

var Gateway = &HandlerStruct{

	ProxyMap: make(map[string]map[int64]*net.ProxyClientStruct),
}

type HandlerStruct struct {
	ProxyMutex sync.Mutex
	//代理地址 -》 客户端连接
	//路由的服务组(LoginProxy,GameProxy)对应的channel  客户端用户id 对应的客户端通道
	ProxyMap map[string]map[int64]*net.ProxyClientStruct
	//登录服务端
	LoginProxy string
	//游戏服务端
	GameProxy string
}

func (handlerBean *HandlerStruct) All(request *net.WsMsgRequestStruct, response *net.WsMsgResponseStruct) {

	name := request.Body.Name
	proxyStr := ""
	if isAccount(name) {
		proxyStr = handlerBean.LoginProxy
	}

	client := net.CreateNewProxyClient(proxyStr)
	client.
		log.Println("路由 * ing")
}
func isAccount(name string) bool {
	if strings.HasPrefix(name, "account.") {
		return true
	} else {
		return false
	}
}
