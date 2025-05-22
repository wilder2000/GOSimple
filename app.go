package main

import (
	"flag"
	"fmt"

	"github.com/wilder2000/GOSimple/http"
)

func main() {
	fmt.Println("GOSimple")
	installCmd := flag.String("install", "NO", "Init database struct and security data.")
	flag.Parse()
	if *installCmd == "YES" {
		hs := http.CreateHttpServer(":9090")
		hs.Install()
	}

}
