package main

import "github.com/cdle/sillyGirl/core"

func main() {
	go core.RunServer()
	select {}
}
