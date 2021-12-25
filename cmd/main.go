package main

import (
	"github.com/dollarkillerx/FollowmeHealthExamination/internel/config"
	"github.com/dollarkillerx/FollowmeHealthExamination/internel/server"
	"github.com/dollarkillerx/FollowmeHealthExamination/utils"
)

func main() {
	err := config.InitConfig("config.json")
	if err != nil {
		panic(err)
	}

	utils.InitCache()

	ser := server.NewServer()
	ser.Run()
}
