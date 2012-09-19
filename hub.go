package main

import (
	"log"
)

type hub struct {
	connections map[string]*connection
	connect     chan *connection
	disconnect  chan *connection
}

var h = hub{
	connections: make(map[string]*connection),
	connect:     make(chan *connection),
	disconnect:  make(chan *connection),
}

// hub.run keeps track of all active connections
func (h *hub) run() {
	log.Printf("[hub.run] starting")
	for {
		select {
		case conn := <-h.connect:
			log.Printf("[hub.run] connecting id: %s", conn.id)
			h.connections[conn.id] = conn
			log.Printf("[hub.run] connections: %v", h.connections)
		case conn := <-h.disconnect:
			log.Printf("[hub.run] disconnecting id: %s", conn.id)
			if conn.id != "" {
				delete(h.connections, conn.id)
			}
		}
	}
}
