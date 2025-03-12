package controllers

import (
	"ThreeKingdoms/src/constant"
	"ThreeKingdoms/src/net"
	"fmt"
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
	log.Println("路由 * ing")
	name := request.Body.Name
	proxyStr := ""
	if isAccount(name) {
		proxyStr = handlerBean.LoginProxy
	} else {
		proxyStr = handlerBean.GameProxy
	}
	if proxyStr == "" {
		response.Body.Code = constant.ProxyNotInConnect
		return
	}
	if proxyStr == "heartbeat" {
		response.Body.Code = constant.OK
		return
	}
	handlerBean.ProxyMutex.Lock()
	proxyMap := handlerBean.ProxyMap[proxyStr]
	if proxyMap == nil {
		proxyMap = make(map[int64]*net.ProxyClientStruct)
		handlerBean.ProxyMap[proxyStr] = proxyMap
	}
	handlerBean.ProxyMutex.Unlock()
	//获取客户端id
	c, err := request.Connection.GetProperty("cid")
	if err != nil {
		log.Println("cid 获取出错")
	}
	cid := c.(int64)
	proxyClient := proxyMap[cid]
	if proxyClient == nil {
		//首次尝试访问服务
		//网关构建代理端

		fmt.Println("首次尝试访问服务: " + proxyStr)
		proxyClient = net.CreateNewProxyClient(proxyStr)
		err := proxyClient.Connect()
		if err != nil {
			handlerBean.ProxyMutex.Lock()
			delete(handlerBean.ProxyMap[proxyStr], cid)
			handlerBean.ProxyMutex.Unlock()
			response.Body.Code = constant.ProxyConnectError
			return
		}
		handlerBean.ProxyMap[proxyStr][cid] = proxyClient
		proxyClient.Channel.SetProperty("cid", cid)
		proxyClient.Channel.SetProperty("proxy", proxyStr)
		proxyClient.Channel.SetProperty("gateway", request.Connection)
		proxyClient.Channel.OnPush = func(conn *net.WsClientChannelStruct, body *net.ResponseStruct) {

		}

	}
	fmt.Println("消息转发ing")
	proxyRseponse, err := proxyClient.Send(request.Body.Name, request.Body.MsgContent)

	if err != nil {
		mapClient := Gateway.ProxyMap[name]
		proxyClient.Channel.Close()
		delete(mapClient, cid)
		log.Println("消息发送出错,删除通道")
		return
	} else {
		fmt.Println("响应结果", proxyRseponse)
		fmt.Println("消息转发成功")
	}
	if response != nil {
		response.Body.Code = proxyRseponse.Code
		response.Body.MsgContent = proxyRseponse.MsgContent
	} else {
		response.Body.Code = constant.ProxyConnectError
		return
	}

}
func isAccount(name string) bool {
	if strings.HasPrefix(name, "account.") {
		return true
	} else {
		return false
	}
}
