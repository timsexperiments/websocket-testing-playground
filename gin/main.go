package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	r := gin.Default()

	r.Use(static.Serve("/chat", static.LocalFile(publicDir(), true)))

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		ipAddress := c.ClientIP()
		conn.WriteMessage(1, []byte(fmt.Sprintf("%s connected to server.", ipAddress)))
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Printf("%s disconnected from server with code [%d] and text [%s]. ", ipAddress, code, text)
			return nil
		})

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	})

	r.GET("/", func(c *gin.Context) {
		fmt.Printf("c.Writer: %v\n", c.Writer)
		fmt.Fprint(c.Writer, "Sup, World. 😎")
	})

	r.Run("0.0.0.0:5000")
}

func publicDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get the current filename.")
	}
	return filepath.Join(filepath.Dir(filename), "/public")
}
