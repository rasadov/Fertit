package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type StaticHandler struct {
	StaticDir string
}

func NewStaticHandler(staticDir string) *StaticHandler {
	return &StaticHandler{
		StaticDir: staticDir,
	}
}

func (s *StaticHandler) Favicon(c *gin.Context) {
	c.File("./web/static/favicon.ico")
}

func (s *StaticHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Year": "2025",
	})
}
