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

const DEFAULT_ROOM  = "wshandler-room"

type Client struct {
	Conn *websocket.Conn
	Room string
	Req  *http.Request
	Res  *http.ResponseWriter
	ID    string
	mu    sync.Mutex
}

func (client *Client) Add() {
	lock.Lock()
	defer lock.Unlock()
	clients[client] = true
}

func (client *Client) Remove() {
	lock.Lock()
	defer lock.Unlock()
	delete(clients, client)
}

func (client *Client) Send(msg []byte, room interface{}) {
	lock.RLock()
	defer lock.RUnlock()

	if room == nil {
		room = DEFAULT_ROOM
	}
	for c := range clients {
		if c.Room == room {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

}

func Broadcast(msg []byte) {
	lock.RLock()
	defer lock.RUnlock()
	for c := range clients {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.Conn.WriteMessage(websocket.TextMessage, msg)
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
		room = DEFAULT_ROOM
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
	OnDisconnect  func(c *Client)
	OnConnect     func(c *Client)
	OnTextMessage func(c *Client, msg []byte)
}

