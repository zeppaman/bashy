package main

import (
	"bashy/src/bashy"
	"bashy/src/utils"
	"fmt"
	"os"
)

func main() {

	fmt.Println("OS:" + utils.CurrentOS())
	instance := bashy.Bashy{}
	err := instance.Init()
	if err == nil {
		instance.Run(os.Args)
	} else {
		fmt.Println("not an happy ending.")
		fmt.Println(err)
	}
	instance.Destroy()
}
