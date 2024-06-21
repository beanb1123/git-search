package main

import (
	"fmt"
	"github.com/agustin-del-pino/git-search/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
