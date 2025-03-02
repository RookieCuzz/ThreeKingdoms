package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type ServerStruct struct {
	addr   string
	router *routerStruct
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
	wsConnection, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	// 确保在函数结束时关闭连接
	defer func() {
		err := wsConnection.Close()
		if err != nil {
			fmt.Println("关闭 WebSocket 连接时出错:", err)
		}
	}()
	//defer wsConnection.Close()
	//err = wsConnection.WriteMessage(websocket.TextMessage, []byte("你好 golang！"))
	//if err != nil {
	//	return
	//}
	for {
		messageType, p, err := wsConnection.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		if messageType == websocket.TextMessage {

			fmt.Println("接受到的消息为:" + string(p))
			wsConnection.WriteMessage(websocket.TextMessage, []byte("echo "+string(p)))
		}
	}

}

func NewServer(addr string) *ServerStruct {
	return &ServerStruct{
		addr:   addr,
		router: nil,
	}
}

//使用开源的websocket协议升级器
var wsUpgrader = websocket.Upgrader{
	//对浏览器发来的请求 允许其跨域行为
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
