package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var adminFS embed.FS

func adminHandler() gin.HandlerFunc {
	subFs, err := fs.Sub(adminFS, "web/dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(subFs))

	return func(c *gin.Context) {
		path := c.Param("filepath")
		if path == "" {
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		if filepath.Ext(path) != "" {
			c.Request.URL.Path = path
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		c.Request.URL.Path = "/"
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
