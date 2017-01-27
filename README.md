# go-wshandler
Simple Go WebSockets Handler


Example:

```
package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"github.com/mikicaivosevic/go-wshandler"
	"github.com/gorilla/mux"
)


func indexHandler(w http.ResponseWriter, r *http.Request) {

	indexFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(index))
}



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
	r.HandleFunc("/", indexHandler)
	http.ListenAndServe(":3000", r)
}

```
`
