package net

import (
	"log"
	"strings"
)

type GroupStruct struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type RouterStruct struct {
	groups []*GroupStruct
}

func (router *RouterStruct) Run(request *WsMsgRequestStruct, response *WsMsgResponseStruct) {
	//包类型
	//登录业务： account.login
	//account组表示
	//login路由标识
	split := strings.Split(request.Body.Name, ".")
	if len(split) == 2 {
		prefix := split[0]
		name := split[1]

		for _, group := range router.groups {
			//执行对应的处理方法
			if group.prefix == prefix {
				handler, ok := group.handlerMap[name]
				if ok {
					handler(request, response)
				}
			} else if group.prefix == "*" {
				handler, ok := group.handlerMap["*"]
				if ok {
					handler(request, response)
				}
			}
		}
		/////组路由
		//switch split[0] {
		//case "account":
		//	//处理账号组
		//	fmt.Printf("Account %s \n", split[1])
		//case "shop":
		//	fmt.Printf("Shop %s \n", split[1])
		//
		//}
	}
}

// 规定一种函数类型
type HandlerFunc func(request *WsMsgRequestStruct, response *WsMsgResponseStruct)

func (group *GroupStruct) exec(name string, request *WsMsgRequestStruct, response *WsMsgResponseStruct) {
	handlerFunc := group.handlerMap[name]
	if handlerFunc != nil {
		handlerFunc(request, response)
	} else {
		handlerFunc = group.handlerMap["*"]

		if handlerFunc != nil {
			handlerFunc(request, response)
		} else {
			log.Println(" * 路由未定义")
		}

	}
}

func (router *RouterStruct) CreateGroup(prefix string) *GroupStruct {
	group := GroupStruct{prefix: prefix, handlerMap: make(map[string]HandlerFunc)}
	router.groups = append(router.groups, &group)
	return &group
}
func (group *GroupStruct) AddEventHandler(name string, handler HandlerFunc) {
	group.handlerMap[name] = handler
}

func NewRouter() *RouterStruct {
	router := &RouterStruct{}
	return router
}
