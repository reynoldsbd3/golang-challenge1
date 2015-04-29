package main

import (
	"drum"
	"fmt"
	"os"
)

func main() {

	pattern := drum.DecodeFile(os.Args[1])

	fmt.Println(pattern)
}
