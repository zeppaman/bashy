package main

import (
	"bashy/src/bashy"
	logger "bashy/src/logger"
	"bashy/src/utils"
	"fmt"
	"os"
)

func main() {
	greenPrintln := logger.CreateColoredPrintln("green")
	greenPrintln("OS: " + utils.CurrentOS())
	instance := bashy.Bashy{}
	err := instance.Init()
	if err == nil {
		if logger.IsDebug() {
			logger.Log("info", "args: "+fmt.Sprintf("%v", os.Args))
		}
		instance.Run(os.Args)
	} else {
		logger.Log("error", "not an happy ending.")
		logger.Log("error", err.Error())
	}
	instance.Destroy()
}
