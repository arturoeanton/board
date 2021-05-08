package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/arturoeanton/board/pkg/chat"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Board v0.0.1")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.StaticFS("/static", http.Dir(""))
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	router.StaticFile("/", "static/index.html")
	router.GET("/stream", chat.Stream)
	router.GET("/msg", chat.EventMessage)
	router.GET("/refresh", chat.EventRefresh)

	router.GET("/chat", func(c *gin.Context) {
		userid := fmt.Sprint(rand.Int31())
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"userid": userid,
		})
	})

	router.Run(":8080")

}
