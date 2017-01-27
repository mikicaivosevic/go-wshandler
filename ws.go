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
var rooms = make(map[string][]*Client)

const DEFAULT_ROOM  = "wshandler-room"

type Client struct {
	Conn *websocket.Conn
	Room string
	Req  *http.Request
	Res  *http.ResponseWriter
	ID string
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

func (client *Client) JoinRoom(room string) {
	lock.Lock()
	defer lock.Unlock()
	client.Room = room
	rooms[room] = append(rooms[room], client)
}

func (client *Client) LeaveRoom(room string) {
	lock.Lock()
	defer lock.Unlock()
	delete(rooms, room)
}

func (client *Client) Send(msg []byte, room interface{}) {
	lock.RLock()
	defer lock.RUnlock()

	if room == nil {
		room = DEFAULT_ROOM
	}

	for roomName, cl := range rooms {
		for _, roomClient := range cl {
			if room == roomName {
				roomClient.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}

}

func Broadcast(msg []byte) {
	lock.RLock()
	defer lock.RUnlock()
	for c := range clients {
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

