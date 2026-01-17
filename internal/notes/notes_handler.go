package notes

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) CreateNote(c *gin.Context) {
	// *gin.Context :- This context hold your information or it holds each request info

	var req CreateNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid json",
		})
		return
	}

	now := time.Now().UTC()
	note := Note{
		ID: primitive.NewObjectID(),

		Title:     req.Title,
		Content:   req.Content,
		Pinned:    req.Pinned,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := h.repo.Create(c.Request.Context(), note)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create note here!",
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) ListNotes(c *gin.Context) {
	notes, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch all notes",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func (h *Handler) getNoteByID(c *gin.Context) {

	idStr := c.Param("id")

	// from the MongoDB driver, converts a 24-character hexadecimal
	// string (like "61791c74138d41367e52d832") into a 12-byte BSON
	// primitive.ObjectID, which is MongoDB's unique document identifier,
	// returning an ObjectID and an error if the string isn't valid,
	// allowing you to use string IDs to query or work with MongoDB
	// documents in your Go application.
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	note, err := h.repo.GetByID(c.Request.Context(), objID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Note not found for that given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch the note",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (h *Handler) UpdateNoteByID(c *gin.Context) {
	idStr := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid json format",
		})
		return
	}

	updated, err := h.repo.UpdateByID(c.Request.Context(), objID, req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Note not found for that given ID",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch the note",
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}
