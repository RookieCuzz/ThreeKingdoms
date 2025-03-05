package main

import (
	"ThreeKingdoms/src/config"
	"ThreeKingdoms/src/gamedatabase"
)

func main() {

	config.Init()
	gamedatabase.TestDatabase()
}
