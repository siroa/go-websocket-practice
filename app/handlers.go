package app

import (
	"embed"
	"log"
	"text/template"
	"websocket/service"

	"github.com/gin-gonic/gin"
)

//go:embed templates/index.html
var indexTmpl embed.FS

type DefaultHandlers struct {
	service service.Service
}

func NewHandlers(s service.Service) DefaultHandlers {
	return DefaultHandlers{s}
}

func (h DefaultHandlers) WebSocketConnection(c *gin.Context) {
	w, r := c.Writer, c.Request
	id := c.Param("id")
	h.service.CreateChatRoom(id, w, r)
}

func (h DefaultHandlers) index(c *gin.Context) {
	tmpl, err := template.ParseFS(indexTmpl, "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(c.Writer, "")
	if err != nil {
		log.Fatal(err)
	}
}
