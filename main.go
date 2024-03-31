package main

import (
	"log"
	"os"

	"github.com/AlvareYN/auto-updater/cmd"
	"github.com/AlvareYN/auto-updater/internal/updater"
	"github.com/gin-gonic/gin"
)

func main() {
	cmd.ConfigInit()

	log.Println("running on pid:", os.Getpid())

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/check-updates", updater.CheckUpdates)

	r.POST("/download-updates", updater.Update)

	r.POST("/apply-updates", updater.ApplyUpdates)

	r.GET("/version", updater.GetVersion)

	r.Run(":8080")
}
