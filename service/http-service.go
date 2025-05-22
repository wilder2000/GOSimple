package service

import (
	"time"

	"github.com/wilder2000/GOSimple/http"
)

type HttpServerConfig struct {
	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
}

type HttpServer struct {
	Address       string
	Config        HttpServerConfig
	Actions       []http.HttpController
	NoAuthActions []http.HttpController
}

func (hs *HttpServer) AppendActions(ac ...http.HttpController) {
	hs.Actions = append(hs.Actions, ac...)
}
func (hs *HttpServer) AppendNoAuthActions(ac ...http.HttpController) {
	hs.NoAuthActions = append(hs.NoAuthActions, ac...)
}
func (hs *HttpServer) Start() {
	http.StartWebServer(hs.Address, *hs.Config.ReadTimeout, *hs.Config.WriteTimeout, hs.Actions, hs.NoAuthActions)
}

func CreateHttpServer(address string, config ...HttpServerConfig) *HttpServer {
	defaultConfig := HttpServerConfig{
		ReadTimeout:  durationPtr(30 * time.Second),
		WriteTimeout: durationPtr(20 * time.Second),
	}
	hs := &HttpServer{
		Config:        defaultConfig,
		Address:       address,
		Actions:       make([]http.HttpController, 0),
		NoAuthActions: make([]http.HttpController, 0),
	}
	if config != nil && len(config) == 1 {
		if config[0].ReadTimeout != nil {
			hs.Config.ReadTimeout = config[0].ReadTimeout
		}
		if config[0].WriteTimeout != nil {
			hs.Config.WriteTimeout = config[0].WriteTimeout
		}
	}

	return hs
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}
