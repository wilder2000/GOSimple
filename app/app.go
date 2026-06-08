package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wilder2000/GOSimple/config"
	"github.com/wilder2000/GOSimple/http"
)

func Run(adminHandler ...gin.HandlerFunc) {
	fmt.Println("GOSimple")

	args := os.Args[1:]
	for _, arg := range args {
		if strings.EqualFold(arg, "-install=YES") || strings.EqualFold(arg, "-install YES") {
			hs := http.CreateHttpServer(config.AConfig.DataSource.Type)
			hs.Install()
			return
		}
		if strings.EqualFold(arg, "-sync-urls") || strings.EqualFold(arg, "-sync-urls=YES") {
			http.SyncUrlMappings()
			return
		}
	}

	hs := http.CreateHttpServer(config.AConfig.Port)
	if len(adminHandler) > 0 {
		hs.AdminHandler = adminHandler[0]
	}
	hs.Start()
}
