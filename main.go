package main

import (
	"fmt"
	"os"
)

func main() {
	// STEP 1: Load the config, if it exists, or create the file
	file := Config()

	// STEP n: Select the files apart from the build folder, to send to the server
	impFiles := SelectFiles(*file)
	if len(impFiles) == 0 {
		fmt.Println("No Files selected")
		os.Exit(0)
	}

	file.Build.Files = impFiles
	file.save()

}
