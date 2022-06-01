package main

import (
	"bashy/src/bashy"
	"fmt"
	"os"
)

func main() {

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
