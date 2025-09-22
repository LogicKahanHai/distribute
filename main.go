package main

import (
	"fmt"
	"os"
)

func main() {
	out := SelectFiles()
	if len(out) == 0 {
		fmt.Println("No Files selected")
		os.Exit(0)
	}
	fmt.Printf("%v", out)
}
