package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type StaticHandler struct {
}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (s *StaticHandler) Favicon(c *gin.Context) {
	c.File("./web/static/favicon.ico")
}

func (s *StaticHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Year": "2025",
	})
}
