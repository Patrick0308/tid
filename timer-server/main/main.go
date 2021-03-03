package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"timer-server"
)

func main() {
	r := gin.Default()
	r.GET("/getSecond", func(c *gin.Context) {
		c.String(200,fmt.Sprintf("%d", timer_server.GetSecond()))
	})
	timer_server.StartTime()
	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}