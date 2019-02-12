package main

import (
	"flag"
	"fmt"

	"github.com/monochromegane/mole"
)

var flagAll bool

func init() {
	flag.BoolVar(&flagAll, "all", false, "Print all versions.")
}

func main() {
	flag.Parse()
	list, err := mole.Run(flagAll)
	if err != nil {
		panic(err)
	}

	for _, dir := range list {
		fmt.Println(dir)
	}
}
