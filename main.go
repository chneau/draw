package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"

	_ "github.com/chneau/draw/pkg/statik"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	if runtime.GOOS == "windows" {
		gin.DisableConsoleColor()
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	gracefulExit()
}

func ce(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}

func gracefulExit() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		os.Exit(0)
	}()
}

func main() {
	port := "3000"
	fs, err := fs.New()
	ce(err, "fs.New()")
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ws", func(c *gin.Context) {
		// handler := websocket.Handler(EchoServer)
		// handler.ServeHTTP(c.Writer, c.Req)
	})
	r.StaticFS("/draw", fs)
	hostname, err := os.Hostname()
	ce(err, "os.Hostname")
	log.Printf("Listening on http://%[1]s:%[2]s/ , http://localhost:%[2]s/\n", hostname, port)
	r.Run(":" + port)
}
