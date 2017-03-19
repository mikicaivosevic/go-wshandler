# go-wshandler
Simple Go WebSockets Handler


####When a message is sent with the Broadcast method, all clients connected to the namespace receive it, including the sender.

    Broadcast(msg []byte)


#### The Send() function accept room argument that cause the message to be sent to all the clients that are in the given room.

    Send(msg []byte, room interface{})

    client.Send(msg, "uuid-room")



Example:

```
package main

import (
	"net/http"
	"fmt"
	"github.com/mikicaivosevic/go-wshandler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)


func wsHandler(w http.ResponseWriter, r *http.Request) {
	//Set websocket upgrader, allow cross domain requests
	wshandler.SetWebSocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	})

	wshandler.WebSocketHandler(w, r, &wshandler.WebSocketEvent{
		OnConnect: func(client *wshandler.Client) {
			fmt.Println("Connected!!!")
		},
		OnDisconnect: func(client *wshandler.Client) {
			fmt.Println(client.Room)
			fmt.Println("Disconnected!!!")
		},
		OnTextMessage: func(client *wshandler.Client, msg []byte) {
			wshandler.Broadcast(msg)
			fmt.Println("Message: " + string(msg))
		},
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":3000", r)
}

```

