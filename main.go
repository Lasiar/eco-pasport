package main

import (
	"eco-passport-back/web"
	"flag"
	"fmt"
)

var (
	printVer bool
	version  string
)

func init() {
	flag.BoolVar(&printVer, "version", false, "Print version")
	flag.Parse()
}

func main() {
	if printVer {
		fmt.Println(version)
	}
	web.Run()
}
