package hub

import (
	"github.com/gorilla/websocket"
)

type msg struct {
	S [4]int `json:"s"`
	C int    `json:"c"`
}

// Hub is the class that takes care of getting all connections
type Hub struct {
	conns map[*websocket.Conn]interface{}
	read  chan *msg
	cache []*msg
}

func (h *Hub) keepDispatching() {
	go func() {
		for {
			m := <-h.read
			for c := range h.conns {
				if c.WriteJSON(m) != nil {
					delete(h.conns, c)
				}
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
			// l := len(h.cache)
			// if l > 100000 {
			// 	h.cache = h.cache[l-90000 : l]
			// }
			h.read <- m
		}
	}()
}

func (h *Hub) initConn(conn *websocket.Conn) {
	go func() {
		for _, c := range h.cache {
			conn.WriteJSON(c)
		}
	}()

}

// AddConn add a connection to the pool of connected
func (h *Hub) AddConn(conn *websocket.Conn) {
	h.conns[conn] = nil
	h.keepReading(conn)
	h.initConn(conn)
}

// New returns an instance
func New() *Hub {
	h := &Hub{
		conns: map[*websocket.Conn]interface{}{},
		read:  make(chan *msg),
	}
	h.keepDispatching()
	return h
}
