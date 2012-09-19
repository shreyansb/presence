package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"strings"
)

type connection struct {
	ws      *websocket.Conn
	c       chan string
	id      string
	friends []string
}

// one handler per websocket connection
func wsHandler(ws *websocket.Conn) {
	log.Printf("[wsHandler] new websocket connection: %v", ws)
	// create a partial connection
	conn := &connection{
		ws: ws,
		c:  make(chan string),
	}
	log.Printf("conn: %v", conn)
	log.Printf("conn.id: %v", conn.id)

	// start listening for incoming messages
	conn.listen()
}

// receive messages from the websocket and do fun things with them
func (conn *connection) listen() {
	log.Printf("[conn.listen] conn: %v", conn)
	for {
		var message string
		err := websocket.Message.Receive(conn.ws, &message)
		if err != nil {
			log.Printf("[conn.listen] error in Receive: %s", err)
			break
		}
		conn.parseMessage(message)
	}
	conn.disconnect()
}

// save the connection's id and tell the hub that this connection is ready
func (conn *connection) connect(id string, friends []string) {
	log.Printf("[conn.connect] id: %s", id)
	conn.id = id
	conn.friends = friends
	h.connect <- conn
	conn.broadcastPresenceOn()
}

// tell the hub this connection no longer exists
// close the connection's channel to other connections and the connection's websocket
func (conn *connection) disconnect() {
	log.Printf("[conn.disconnect] id: %s", conn.id)
	h.disconnect <- conn
	close(conn.c)
	conn.ws.Close()
	conn.broadcastPresenceOff()
}

// see if the incoming message is the websocket identifying itself, or just sending a message
func (conn *connection) parseMessage(message string) {
	if strings.HasPrefix(message, "id: ") {
		id := strings.Split(message, " ")[1]
		friends := []string{"1", "2", "3", "4"}
		conn.connect(id, friends)
	} else {
		conn.broadcastMessage(message)
	}
}

// send a message to all the user's friends
func (conn *connection) broadcastMessage(message string) {
	log.Printf("[conn.broadcastMessage] message: %s", message)
	for _, id := range conn.friends {
		if id == conn.id {
			continue
		}
		if friend := h.connections[id]; friend != nil {
			log.Printf("[conn.broadcastMessage] sending message to id: %s", id)
			websocket.Message.Send(friend.ws, message)
		}
	}
}

// inform all the friends of the connected user that this user is online
func (conn *connection) broadcastPresenceOn() {
	log.Printf("[conn.broadcastPresenceOn] id: %s", conn.id)
	conn.broadcastMessage(fmt.Sprintf("id %s connected", conn.id))
}

// inform all the friends of the connected user that this user is no longer online
func (conn *connection) broadcastPresenceOff() {
	log.Printf("[conn.broadcastPresenceOff] id: %s", conn.id)
	conn.broadcastMessage(fmt.Sprintf("id %s disconnected", conn.id))
}
