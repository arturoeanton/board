package chat

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/gin-gonic/gin"
)

type Events struct {
	Messages map[string]chan string
	Refresh  map[string]chan string
}

var events *Events
var messagesStorage []string

func Stream(c *gin.Context) {
	if events == nil {
		events = &Events{
			Messages: make(map[string]chan string),
			Refresh:  make(map[string]chan string),
		}
	}
	id := fmt.Sprint(rand.Int31())

	message := make(chan string)
	refresh := make(chan string)
	defer close(message)
	defer close(refresh)

	events.Messages[id] = message
	events.Refresh[id] = refresh
	defer delete(events.Messages, id)
	defer delete(events.Refresh, id)

	c.Stream(func(w io.Writer) bool {
		select {
		case <-refresh:
			c.SSEvent("refresh", messagesStorage)
			return true
		case msg := <-message:
			c.SSEvent("message", msg)
			return true
		}
	})
}

func EventMessage(c *gin.Context) {
	msg := c.Query("msg")
	messagesStorage = append(messagesStorage, msg)
	for _, l := range events.Messages {
		l <- msg
	}
	c.JSON(200, gin.H{
		"data":  msg,
		"staus": "ok",
	})
}

func EventRefresh(c *gin.Context) {
	for _, l := range events.Refresh {
		l <- "refresh"
	}
	c.JSON(200, gin.H{
		"data":  "refresh",
		"staus": "ok",
	})
}
