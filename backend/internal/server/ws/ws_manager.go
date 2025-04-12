// internal/server/ws/ws_manager.go

package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
	send chan Event
}

type symbolHub struct {
	symbol  string
	clients map[*client]bool
	lock    sync.RWMutex
}

var hubs = make(map[string]*symbolHub)
var hubsLock sync.RWMutex

// Publish envia evento para todos os clientes do símbolo
func Publish(symbol string, event Event) {
	hub := getHub(symbol)
	hub.lock.RLock()
	defer hub.lock.RUnlock()

	for c := range hub.clients {
		c.send <- event
	}
}

// getHub retorna (ou cria) o hub de um símbolo
func getHub(symbol string) *symbolHub {
	hubsLock.Lock()
	defer hubsLock.Unlock()

	if hub, ok := hubs[symbol]; ok {
		return hub
	}

	newHub := &symbolHub{
		symbol:  symbol,
		clients: make(map[*client]bool),
	}
	hubs[symbol] = newHub
	return newHub
}

// AddClient adiciona um cliente WebSocket ao hub do símbolo
func AddClient(symbol string, conn *websocket.Conn) {
	hub := getHub(symbol)

	c := &client{
		conn: conn,
		send: make(chan Event, 10),
	}

	hub.lock.Lock()
	hub.clients[c] = true
	hub.lock.Unlock()

	go c.writeLoop()
}

func (c *client) writeLoop() {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
