package main

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type msg struct {
	S [4]uint16 `json:"s"`
	C uint8     `json:"c"` // color
	W uint8     `json:"w"` // width
}

// Hub is the class that takes care of getting all connections
type Hub struct {
	conns              map[*websocket.Conn]*sync.Mutex
	cache              []*msg
	CacheSize          int
	LatestModification time.Time
	room               string
}

func (h *Hub) Kill() {
	for c, mu := range h.conns {
		func() {
			mu.Lock()
			defer mu.Unlock()
			_ = c.Close()
		}()
	}
}

func (h *Hub) broadcast(conn *websocket.Conn) {
	go func() {
		defer log.Println(h.room, "Broadcast done")
		for {
			m := &msg{}
			err := conn.ReadJSON(m)
			if err != nil {
				return
			}
			h.cache = append(h.cache, m)
			l := len(h.cache)
			if l > h.CacheSize {
				h.cache = h.cache[l-int(float64(h.CacheSize)*0.8) : l]
				log.Println(h.room, "Cache resized:", len(h.cache))
			}
			for c, mu := range h.conns {
				mu.Lock()
				if c.WriteJSON([]*msg{m}) != nil {
					delete(h.conns, c)
					log.Println(h.room, "ws:", len(h.conns))
				}
				mu.Unlock()
			}
			h.LatestModification = time.Now()
		}
	}()
}

func (h *Hub) initConn(conn *websocket.Conn) {
	go func() {
		mu := h.conns[conn]
		log.Println(h.room, "ws:", len(h.conns))
		mu.Lock()
		_ = conn.WriteJSON(h.cache)
		mu.Unlock()
		h.LatestModification = time.Now()
	}()

}

// AddConn add a connection to the pool of connected
func (h *Hub) AddConn(conn *websocket.Conn) {
	h.conns[conn] = &sync.Mutex{}
	h.broadcast(conn)
	h.initConn(conn)
}

// NewHub returns an instance
func NewHub(room string) *Hub {
	h := &Hub{
		conns:              map[*websocket.Conn]*sync.Mutex{},
		cache:              []*msg{},
		room:               room,
		CacheSize:          5000,
		LatestModification: time.Now(),
	}
	return h
}
