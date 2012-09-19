package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

func main() {
	go h.run()

	http.Handle("/ws", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("[main] ListenAndServe error: ", err)
	}
}
