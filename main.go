package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mattn/go-colorable"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	_ "embed"
)

//go:embed index.html
var indexHtml []byte

func init() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = colorable.NewColorableStdout()
	log.SetFlags(log.Ltime)
	log.SetPrefix("[draw] ")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	hubs := map[string]*Hub{}
	mtx := &sync.Mutex{}
	period := time.Second * 10
	inactive := time.Minute * 2
	go func() {
		every := time.NewTicker(period)
		for range every.C {
			func() {
				mtx.Lock()
				defer mtx.Unlock()
				for k, v := range hubs {
					if time.Since(v.LatestModification) > inactive {
						delete(hubs, k)
						v.Kill()
						log.Println(k, "Deleted room")
					}
				}
			}()
		}
	}()
	log.Println("Will check every", period, "if room inactive for more than", inactive, "and will clean it")
	r := gin.New()
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	r.GET("/ws/*room", func(c *gin.Context) {
		mtx.Lock()
		defer mtx.Unlock()
		room := c.Param("room")
		if hubs[room] == nil {
			hubs[room] = NewHub(room)
			log.Println(room, "Created room:")
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())
		if err != nil {
			log.Println("ERROR", err)
			return
		}
		hubs[room].AddConn(conn)
	})
	r.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHtml)
	})
	hostname, _ := os.Hostname()
	log.Printf("Listening on (hostname) http://%[1]s:%[2]s/", hostname, port)
	err := r.Run(":" + port)
	if err != nil {
		log.Panicln(err)
	}
}
