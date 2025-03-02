package net

type groupStruct struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

type routerStruct struct {
	router *[]routerStruct
}

type HandlerFunc func()
