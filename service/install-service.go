package service

import (
	"flag"

	"github.com/wilder2000/GOSimple/glog"
	"github.com/wilder2000/GOSimple/http"
)

func PrepareInstall() {
	installCmd := flag.String("install", "NO", "Init database struct and security data.")
	flag.Parse()
	if *installCmd == "YES" {
		hs := http.CreateHttpServer(":9090")
		hs.Install()
	} else {
		glog.Logger.Info("not found install arg para")
	}
}
