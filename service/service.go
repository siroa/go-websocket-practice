package service

import (
	"net/http"
	"websocket/domain"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// primary port
type Service interface {
	CreateChatRoom(string, gin.ResponseWriter, *http.Request)
	PushMessage()
}

type DefaultService struct {
	repo domain.Repository
	cr   *domain.Chatroom
}

func (s DefaultService) CreateChatRoom(id string, w gin.ResponseWriter, r *http.Request) {
	s.cr.RoomId = id
	websocket.Handler(s.cr.WebsocketHandler).ServeHTTP(w, r)
}

func (s DefaultService) PushMessage() {

}

// factory function
func NewService(repo domain.Repository, cr *domain.Chatroom) DefaultService {
	return DefaultService{repo, cr}
}
