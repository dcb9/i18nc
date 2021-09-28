package main

import (
	"fmt"
	"os"

	"github.com/dcb9/i18nc"
)

func main() {
	pkg := os.Args[1]
	destFile := os.Args[2]
	localeFile := os.Args[3]

	dest, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	if err := i18nc.FromFile(localeFile, pkg, dest); err != nil {
		panic(err)
	}

	fmt.Println("Generated successfully")
}
