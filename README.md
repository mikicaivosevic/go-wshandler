# go-wshandler
Simple Go WebSockets Handler


Example:

```
package main

import (
	"net/http"
	"fmt"
	"github.com/mikicaivosevic/go-wshandler"
	"github.com/gorilla/mux"
)


func wsHandler(w http.ResponseWriter, r *http.Request) {

	wshandler.WebSocketHandler(w, r, &wshandler.WebSocketEvent{
		OnConnect: func(client *wshandler.Client) {
			client.Add()
			fmt.Println("Connected!!!")
		},
		OnDisconnect: func(client *wshandler.Client) {
			client.Remove()
			fmt.Println(client.Room)
			fmt.Println("Disconnected!!!")
		},
		OnTextMessage: func(client *wshandler.Client, msg []byte) {
			client.Send(msg, nil, true)
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

