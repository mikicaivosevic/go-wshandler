package wshandler

import (
	"net/http"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

var upgrader websocket.Upgrader
var clients = make(map[*Client]bool)
var lock = sync.RWMutex{}

type Client struct {
	Conn *websocket.Conn
	Room string
	Req  *http.Request
	Res  *http.ResponseWriter
}

func (client *Client) Add() {
	lock.Lock()
	defer lock.Unlock()
	clients[client]=true
}

func (client *Client) Remove() {
	lock.Lock()
	defer lock.Unlock()
	delete(clients, client)
}

func (client *Client) JoinRoom(room string) {
	client.Room = room
}

func (client *Client) Send(msg []byte, room interface{}, broadcast interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	if broadcast == true {
		for c := range clients {
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
	if broadcast == false {
		client.Conn.WriteMessage(websocket.TextMessage, msg)
	}

	if broadcast == nil {
		for c := range clients {
			if c.Room == room {
				c.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}

}


func WebSocketHandler(w http.ResponseWriter, r *http.Request, OnEvent *WebSocketEvent) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	room := r.URL.Query().Get("room")
	if room == "" {
		room = "wshandler-room"
	}
	client := Client{
		Conn: conn,
		Room: room,
		Req: r,
		Res: &w,
	}

	OnEvent.OnConnect(&client)
	client.Add()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			OnEvent.OnDisconnect(&client)
			client.Remove()
			return
		}
		OnEvent.OnTextMessage(&client, msg)
	}
}

type WebSocketEvent struct {
	OnDisconnect func(c *Client)
	OnConnect func(c *Client)
	OnTextMessage func(c *Client, msg []byte)
}

