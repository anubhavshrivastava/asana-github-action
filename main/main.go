package main

import (
	"fmt"
	"os"
)

func main() {
	name := os.Args[1]
	time := os.Args[2]
	fmt.Printf("::set-output name=time::%s\n", time)
	fmt.Printf("::debug Hey Man you said :%s\n", name)
}
