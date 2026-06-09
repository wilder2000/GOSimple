package http

import (
	"embed"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed webdist
var adminFS embed.FS

var adminHandler gin.HandlerFunc

func init() {
	subFs, err := fs.Sub(adminFS, "webdist")
	if err != nil {
		panic(err)
	}
	adminHandler = createAdminHandler(subFs)
}

func createAdminHandler(subFs fs.FS) gin.HandlerFunc {
	fileServer := http.FileServer(http.FS(subFs))
	return func(c *gin.Context) {
		path := c.Param("filepath")
		if path == "" || path == "/" || filepath.Ext(path) == "" {
			c.Request.URL.Path = "/"
		} else {
			c.Request.URL.Path = path
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
