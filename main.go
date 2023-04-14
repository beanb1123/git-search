package main

import (
	"fmt"
	"github.com/agustin-del-pino/gss/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
