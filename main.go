package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/mattn/go-colorable"

	"github.com/chneau/draw/pkg/hub"

	_ "github.com/chneau/draw/pkg/statik"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rakyll/statik/fs"
)

func init() {
	gracefulExit()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = colorable.NewColorableStdout()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fs, _ := fs.New()
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
	r.StaticFS("/draw", fs)
	printIP(port)
	r.Run(":" + port)
}

func printIP(port string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.To4() == nil {
				continue
			}
			log.Printf("Listening on (%s) http://%s:%s/", i.Name, ip, port)
		}
	}
	hostname, _ := os.Hostname()
	log.Printf("Listening on (hostname) http://%[1]s:%[2]s/", hostname, port)
}
