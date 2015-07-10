package main

import (
	"fmt"
	"log"

	"github.com/bitrise-io/go-pathutil/pathutil"
)

func main() {
	currWorkAbsPth, err := pathutil.CurrentWorkingDirectoryAbsolutePath()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("CurrentWorkingDirectoryAbsolutePath:", currWorkAbsPth)
}
