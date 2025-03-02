package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/net"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	config.Init()
	port := config.Config.LoginServer.Port
	host := config.Config.LoginServer.Host
	wg.Add(1)
	net.NewServer(host + ":" + strconv.Itoa(port)).Start()
}
