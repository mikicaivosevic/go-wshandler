# go-wshandler
Simple Go WebSockets Handler


####When a message is sent with the broadcast argument true, all clients connected to the namespace receive it, including the sender.

    Send(msg []byte, room interface{}, broadcast interface{})

    client.Send(msg, "room", true)


#### The Send() function accept room argument that cause the message to be sent to all the clients that are in the given room if broadcast argument is nil.

    Send(msg []byte, room interface{}, broadcast interface{})

    client.Send(msg, "uuid-room", nil)



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
			fmt.Println("Connected!!!")
		},
		OnDisconnect: func(client *wshandler.Client) {
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

