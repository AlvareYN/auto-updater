package main

import (
	"github.com/AlvareYN/auto-updater/cmd"
	"github.com/gin-gonic/gin"
)

func checkUpdatesHandler(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Checking for updates",
	})
}

func Main() {
	cmd.ConfigInit()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/check-updates", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Checking for updates",
		})
	})

	r.POST("/update")

	r.Run(":8080")
}
