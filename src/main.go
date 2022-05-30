package main

import (
	"bashy/src/bashy"
	"fmt"
	"os"
)

func main() {

	bashy := bashy.Bashy{}
	app, err := bashy.Init()
	if err == nil {
		app.Run(os.Args)
	} else {
		fmt.Println("not an happy ending.")
		fmt.Println(err)
	}
}
