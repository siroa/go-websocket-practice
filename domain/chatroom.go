package domain

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

/*
チャットルーム・モデル
*/
type Chatroom struct {
	forward chan *Message //boardcast    client = subscriptio
	join    chan *client  //regsiter
	leave   chan *client  //Unregister
	clients map[string]map[*client]bool
	RoomId  string
}

func (c *Chatroom) connection(gc *gin.Context) {
	w, r := gc.Writer, gc.Request
	id := gc.Param("id")
	c.RoomId = id
	websocket.Handler(c.WebsocketHandler).ServeHTTP(w, r)
}

/*
Chatroomをhttp.handleに適合させる。
ここでは以下のことを実装する。
	・websocketの開設
	・clientの生成
*/
func (c *Chatroom) WebsocketHandler(ws *websocket.Conn) {
	defer ws.Close()
	/* クライアントの生成 */
	client := &client{
		socket:  ws,
		send:    make(chan *Message),
		room:    c,
		roomKey: c.RoomId,
	}

	// チャットルームのjoinチャネルにアクセスし、クライアントを入室させる。最後には必ず退室させる。
	c.join <- client
	fmt.Println("入室")
	defer func() {
		c.leave <- client
		client.exit()
		//c.leave <- client
		//ws.Close()
	}()

	go client.write()
	client.read()
}

// チャットルームの作成
func NewRoom() *Chatroom {
	return &Chatroom{
		forward: make(chan *Message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[string]map[*client]bool),
	}
}

func (c *Chatroom) Run() {

	// チャットルームは無限にルームで起動させる
	for {
		// チャンネルの動きを監視し、処理を決定する
		select {

		// joinチャネルに動きがあった場合（クライアントの入室)
		case lclient := <-c.join:
			conn := c.clients[lclient.roomKey]
			if conn == nil {
				conn = make(map[*client]bool)
				c.clients[lclient.roomKey] = conn
			}
			c.clients[lclient.roomKey][lclient] = true
			fmt.Printf(lclient.roomKey+"に"+"クライアントが入室しました。現在 %x 人のクライアントが存在します。\n",
				len(c.clients[lclient.roomKey]))
		// leaveチャネルに動きがあった場合（クライアントの退室）
		case lclient := <-c.leave:
			delete(c.clients[lclient.roomKey], lclient)
			fmt.Printf(lclient.roomKey+"の"+"クライアントが退室しました。現在 %x 人のクライアントが存在します。\n",
				len(c.clients[lclient.roomKey]))

		// forwardチャネルに動きがあった場合（メッセージの受信）
		case msg := <-c.forward:
			fmt.Println("メッセージを受信しました")
			//存在するクライアント全てに対してメッセージを送信する
			for target := range c.clients[msg.roomKey] {
				select {
				case target.send <- msg:
					fmt.Println("メッセージの送信に成功しました")
				default:
					fmt.Println("メッセージの送信に失敗しました")
					delete(c.clients[msg.roomKey], target)
				}
			}
		}
	}
}
