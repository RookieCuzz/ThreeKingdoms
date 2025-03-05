package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type ServerStruct struct {
	addr   string
	Router *RouterStruct
}

func (server *ServerStruct) Start() {
	//添加url和对应处理函数
	http.HandleFunc("/", server.wsHandler)

	fmt.Printf("服务器成功启动！ %s\n", server.addr)
	//开启端口并监听
	err := http.ListenAndServe(server.addr, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}
func (server *ServerStruct) wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("who is coming？")

	//升级
	wsConnection, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	//构建链接通道
	wsServer := NewWsServer(wsConnection)
	wsServer.Router(server.Router)
	//开启循环
	wsServer.Start()
	wsServer.Handshake()
	fmt.Printf("握手没?")
	// 确保在函数结束时关闭连接
	defer func() {
		err := wsConnection.Close()
		if err != nil {
			fmt.Println("关闭 WebSocket 连接时出错:", err)
		}
	}()
	for {
	}

}

func NewServer(addr string) *ServerStruct {
	return &ServerStruct{
		addr:   addr,
		Router: nil,
	}
}

//使用开源的websocket协议升级器
var wsUpgrade = websocket.Upgrader{
	//对浏览器发来的请求 允许其跨域行为
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
