package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/spotify/oauth", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "hello world!"})
	})

	r.Run()
}
