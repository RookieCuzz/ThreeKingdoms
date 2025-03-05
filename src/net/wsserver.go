package net

import (
	"ThreeKingdoms/src/services/proto"
	"ThreeKingdoms/src/utils"
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

//go的结构体可以理解为在堆空间(或者栈空间)的的数据(是有状态的,不同时刻内容可能是不一样的)
//而方法是无状态的
//所以go函数的参数如果是想要使用传入结构体的数据 应该传入指针
//如果只是想要使用这个结构体实现的某些方法可以不传(因为方法是没有状态的)

// 其实可以理解为对 用户-服务器的连接的进一步封装
type wsServerChannelStruct struct {
	wsConnection *websocket.Conn
	//消息的路由
	router *RouterStruct

	//FIFO写队列 用通道是因为方便不同协程写入读取数据
	outChannelBuffer chan *WsMsgResponseStruct
	Seq              int64
	properties       map[string]interface{}
	//写锁
	propertyLock sync.RWMutex
}

// 客户端与服务端建立的websocket通道需要先进行握手
func (wsServerChannelBean *wsServerChannelStruct) Handshake() {
	//尝试握手

	secrectKey, err := wsServerChannelBean.GetProperty("secretKey")
	if err == nil && secrectKey != nil {
		secrectKey = secrectKey.(string)
	} else {
		//初始化秘钥
		secrectKey = utils.RandSeq(16)
		wsServerChannelBean.SetProperty("secretKey", secrectKey)
	}
	secrectKey2, err := wsServerChannelBean.GetProperty("secretKey")
	fmt.Printf("秘钥是:%s", secrectKey2)
	handSocketPack := &proto.LoginHandshakeServerPacketStruct{
		Key: secrectKey.(string),
	}
	body := &ResponseStruct{
		Name:       "handshake",
		MsgContent: handSocketPack,
	}
	//json化
	jsonBody, _ := json.Marshal(body)
	//压缩
	zipContent, err := utils.Zip(jsonBody)
	if err != nil {
		log.Fatal(err)
	}
	wsServerChannelBean.wsConnection.WriteMessage(websocket.BinaryMessage, zipContent)
}

func (wsServerBean wsServerChannelStruct) SetProperty(key string, value interface{}) {
	//添加属性
	//上隔离写 同时不共享读
	wsServerBean.propertyLock.Lock()

	//保证方法结束回收锁
	defer wsServerBean.propertyLock.Unlock()
	//设置属性
	wsServerBean.properties[key] = value

}

func (wsServerBean wsServerChannelStruct) GetProperty(key string) (interface{}, error) {
	//获取属性 上读锁 共享读
	wsServerBean.propertyLock.RLock()
	defer wsServerBean.propertyLock.RUnlock()
	//返回数据
	return wsServerBean.properties[key], nil

}

func (wsServerBean wsServerChannelStruct) RemoveProperty(key string) {
	//删除数据
	wsServerBean.propertyLock.Lock()
	defer wsServerBean.propertyLock.Unlock()

	//1.
	delete(wsServerBean.properties, key)

}

func (wsServerBean wsServerChannelStruct) Addr() string {
	//客户端的ip:端口(temp)
	return wsServerBean.wsConnection.RemoteAddr().String()
}

func (wsServerBean wsServerChannelStruct) Push(name string, data interface{}) {

	rsp := &WsMsgResponseStruct{
		Body: &ResponseStruct{
			Name:       name,
			MsgContent: data,
			Seq:        0,
		},
	}

	//写入缓存
	wsServerBean.outChannelBuffer <- rsp
}

func (wsServerStruct *wsServerChannelStruct) Router(router *RouterStruct) {
	wsServerStruct.router = router
}

func NewWsServer(connetion *websocket.Conn) *wsServerChannelStruct {

	return &wsServerChannelStruct{
		wsConnection:     connetion,
		outChannelBuffer: make(chan *WsMsgResponseStruct, 1000),
		properties:       make(map[string]interface{}),
		Seq:              0,
	}
}

func (wsServerBean wsServerChannelStruct) Start() {

	//启动读写循环协程
	go flushMsgLoop(&wsServerBean)
	go readMsgLoop(&wsServerBean)
}
func (wsServerBean wsServerChannelStruct) Stop() {}

// 将服务端中的缓冲区消息通过socket发送到客户端
func flushMsgLoop(wsServerBean *wsServerChannelStruct) {
	for {
		select {
		//尝试 操作1 若成功则执行case下语句并跳出
		//这里是尝试从通道读取数据

		//写之前应该加密压缩
		case msg := <-wsServerBean.outChannelBuffer:

			//将msg json化
			data, err2 := json.Marshal(msg.Body)
			if err2 != nil {
				log.Println(err2)
			}

			//加压 加密
			data = processServerMessage(data, wsServerBean)
			wsServerBean.wsConnection.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}

func processServerMessage(data []byte, wsServerBean *wsServerChannelStruct) []byte {

	//加密
	//1获取秘钥
	secretKey, err := wsServerBean.GetProperty("secretKey")
	if err != nil && secretKey == nil {
		log.Fatal(err)
	}
	key := secretKey.(string)
	encrypt, err := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	if err != nil {
		log.Fatal(err)
	}
	//2压缩
	zip, err2 := utils.Zip(encrypt)
	if err2 != nil {
		log.Println(err2)
	}
	return zip
}

func processClientMessage(data []byte, wsServerBean *wsServerChannelStruct) []byte {
	var out_data []byte
	//收到消息 解析消息,前端发来的消息为json消息
	unzipData, err2 := utils.UnZip(data)
	if err2 != nil {
		log.Println("解析请求发生错误", err2)
	}
	//进行解密
	secretKey, err := wsServerBean.GetProperty("secretKey")
	if err == nil {
		//转为字符串
		key := secretKey.(string)

		decrypt, err := utils.AesCBCDecrypt(unzipData, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			//说明 客户端和服务端秘钥不一致
			//服务端再次和客户端进行握手
			wsServerBean.Handshake()

		} else {
			//获取解密数据
			out_data = decrypt
		}
	}

	return out_data

}

// 写入消息
func readMsgLoop(wsServerBean *wsServerChannelStruct) {

	//异常则关闭连接
	defer func() {
		wsServerBean.Close()
		log.Println(wsServerBean.Addr() + " 连接已经关闭")
	}()
	for {
		messageType, data, err := wsServerBean.wsConnection.ReadMessage()
		//判断是否握手

		if err != nil {
			log.Println("收消息发生错误", err)
			break
		}
		if messageType == websocket.TextMessage {
			log.Println(string(data))
		}

		//unzip and crypt
		data = processClientMessage(data, wsServerBean)

		body := &RequestStruct{}
		err = json.Unmarshal(data, body)
		fmt.Println(body.Name)
		if err != nil {
			log.Fatalf("解析json请求发生错误,客户端请检查格式", err)
		}
		//1.将玩家发来的消息进行解密 转json并封装为request
		request := &WsMsgRequestStruct{
			Connection: wsServerBean,
			Body:       body,
		}
		response := &WsMsgResponseStruct{
			Body: &ResponseStruct{
				Name: body.Name,
				Seq:  body.Seq,
				Code: 0,
			},
		}
		//进行解析
		wsServerBean.router.Run(request, response)
		fmt.Printf("发送@@@", response.Body)
		//将请求送入缓冲区
		wsServerBean.outChannelBuffer <- response

		//将request派发给对应的业务线
		//if messageType == websocket.TextMessage {
		//	//得到路由后的消息
		//	//路由后处理消息
		//	fmt.Println("处理消息")
		//	response := &WsMsgResponseStruct{
		//		Body: &ResponseStruct{
		//			Name:       "CPDD",
		//			MsgContent: "这是对 " + string(p) + "的回复",
		//		},
		//	}
		//	wsServerBean.outChannelBuffer <- response
		//}

	}

}
func (wsServerBean *wsServerChannelStruct) Close() {
	wsServerBean.wsConnection.Close()
}
