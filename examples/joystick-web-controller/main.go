package main

import (
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	service := &HTTPService{}
	err := service.Init()
	if err != nil {
		panic(err)
	}
	err = service.Run()
	if err != nil {
		panic(err)
	}

}
