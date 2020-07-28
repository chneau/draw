package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mattn/go-colorable"

	"github.com/chneau/draw/pkg/hub"
	"github.com/chneau/draw/pkg/static"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func init() {
	quit := make(chan os.Signal)
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
	hub := hub.New()
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ws", func(c *gin.Context) {
		conn, _ := websocket.Upgrade(c.Writer, c.Request, c.Writer.Header(), 1024, 1024)
		hub.AddConn(conn)
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(307, "/draw")
	})
	r.StaticFS("/draw", static.AssetFile())
	hostname, _ := os.Hostname()
	log.Printf("Listening on (hostname) http://%[1]s:%[2]s/", hostname, port)
	err := r.Run(":" + port)
	if err != nil {
		log.Panicln(err)
	}
}
