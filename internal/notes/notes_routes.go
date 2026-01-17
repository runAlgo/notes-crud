package notes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	// create repo and handler once at startup
	repo := NewRepo(db)
	h := NewHandler(repo)

	notesGroup := r.Group("/notes")
	{
		notesGroup.POST("", h.CreateNote)
		notesGroup.GET("", h.ListNotes)
		notesGroup.GET("/:id", h.getNoteByID)
		notesGroup.PUT("/:id", h.UpdateNoteByID)
	}

}