package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(app *App) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		if err := pingWithTimeout(app.DB, 2*time.Second); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/users", app.createUser)
	r.GET("/users", app.listUsers)
	r.GET("/users/:id", app.getUser)
	r.PUT("/users/:id", app.updateUser)
	r.DELETE("/users/:id", app.deleteUser)

	return r
}
