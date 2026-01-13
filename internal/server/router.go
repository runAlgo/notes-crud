package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/runAlgo/notes-api/internal/notes"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(database *mongo.Database) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"ok": true,
			"status": "healthy",
		})
	})

	notes.RegisterRoutes(r, database)
	return r
}