package main

import (
	"fmt"

	"github.com/wilder2000/GOSimple/http"
)

func main() {
	fmt.Println("GOSimple")

	hs := http.CreateHttpServer(":9090")
	hs.Start()

}
