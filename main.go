package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("./balajis [directory]")
		return
	}
	PrintEdges(os.Args[1])
}
