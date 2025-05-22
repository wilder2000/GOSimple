package main

import (
	"fmt"

	"github.com/wilder2000/GOSimple/service"
)

func main() {
	fmt.Println("GOSimple")

	hs := service.CreateHttpServer(":9090")
	hs.Start()

}
