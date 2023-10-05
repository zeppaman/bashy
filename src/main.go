package main

import (
	"bashy/src/bashy"
	"bashy/src/logger"
	"bashy/src/utils"
	"os"
)

func main() {
	logger.GreenPrintln("OS: " + utils.CurrentOS())
	logger.BLogInfo("Bashy is starting...")

	instance := bashy.Bashy{}
	err := instance.Init()
	if err == nil {
		instance.Run(os.Args)
	} else {
		logger.BLogError("not an happy ending.")
		logger.BLogError(err)
	}
	instance.Destroy()
	defer logger.CloseLogger()
}
