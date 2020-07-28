package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type msg struct {
	S [4]uint16 `json:"s"`
	C uint8     `json:"c"` // color
	W uint8     `json:"w"` // width
}

// Hub is the class that takes care of getting all connections
type Hub struct {
	conns     map[*websocket.Conn]*sync.Mutex
	read      chan *msg
	cache     []*msg
	CacheSize int
}

func (h *Hub) keepDispatching() {
	go func() {
		for {
			m := <-h.read
			for c, mu := range h.conns {
				mu.Lock()
				if c.WriteJSON([]*msg{m}) != nil {
					delete(h.conns, c)
					log.Println("ws:", len(h.conns))
				}
				mu.Unlock()
			}
		}
	}()
}

func (h *Hub) keepReading(conn *websocket.Conn) {
	go func() {
		for {
			m := &msg{}
			err := conn.ReadJSON(m)
			if err != nil {
				break
			}
			h.cache = append(h.cache, m)
			l := len(h.cache)
			if l > h.CacheSize {
				h.cache = h.cache[l-int(float64(h.CacheSize)*0.8) : l]
				log.Println("Cache resized:", len(h.cache))
			}
			h.read <- m
		}
	}()
}

func (h *Hub) initConn(conn *websocket.Conn) {
	go func() {
		mu := h.conns[conn]
		log.Println("ws:", len(h.conns))
		mu.Lock()
		_ = conn.WriteJSON(h.cache)
		mu.Unlock()
	}()

}

// AddConn add a connection to the pool of connected
func (h *Hub) AddConn(conn *websocket.Conn) {
	h.conns[conn] = &sync.Mutex{}
	h.keepReading(conn)
	h.initConn(conn)
}

// New returns an instance
func New() *Hub {
	h := &Hub{
		conns:     map[*websocket.Conn]*sync.Mutex{},
		read:      make(chan *msg),
		CacheSize: 5000,
	}
	h.keepDispatching()
	return h
}
