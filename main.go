package main

import (
	"ginchat/Router"
	"ginchat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()
	r := Router.App()
	r.Run(":8080")
}
