package main

import (
	"fmt"
	"os"
)

/*
run


*/

const (
	programName = "procman"
)

func main() {

	app := Setup()
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("%v error: %v", programName, err.Error())
		println(errOut)
		os.Exit(1)
	}

}
