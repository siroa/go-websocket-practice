package domain

import (
	"encoding/json"
	"log"
	"websocket/utils/errs"

	"golang.org/x/net/websocket"
)

type Message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	UserId  string
	Time    string
	roomKey string
}

type client struct {
	socket   *websocket.Conn // websocketへのコネクション
	send     chan *Message
	room     *Chatroom // 所属するチャットルーム
	roomKey  string    // ルームのID
	username string    // クライアントの名前
}

// secondary port
type Repository interface {
	saveMessage(Message) *errs.AppError
}

func (c *client) exit() {
	var msg *Message
	msg = &Message{}
	msg.Name = c.username
	msg.roomKey = c.roomKey

	c.room.forward <- msg
}

// websocketに書き出されたメッセージを読み込む
func (c *client) read() {
	defer c.socket.Close()
	for {
		// メッセージを受信する
		var msg Message
		var receiveMsg []byte
		if err := websocket.Message.Receive(c.socket, &receiveMsg); err == nil {
			if err = json.Unmarshal(receiveMsg, &msg); err != nil {
				log.Fatalln(err)
			}
			c.username = msg.Name
			msg.roomKey = c.roomKey
			c.room.forward <- &msg
		} else {
			break
		}
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		sendMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatalln(err)
		}
		err = websocket.Message.Send(c.socket, sendMsg)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
