package notes

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repo -> data access layer

type Repo struct {
	coll *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		coll: db.Collection("notes"),
	}
}

func (r *Repo) Create(ctx context.Context, note Note) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	_, err := r.coll.InsertOne(opCtx, note)
	if err != nil {
		return Note{}, fmt.Errorf("Insert note failed")
	}
	return note, nil
}

func (r *Repo) List(ctx context.Context) ([]Note, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filer := bson.M{} // match all docs

	// Find returns a cursor (like an iterator) -> over matching elements
	cursor, err := r.coll.Find(ctx, filer)

	if err != nil {
		return nil, fmt.Errorf("find notes failed: %w", err)
	}

	// cursor must be closed after use
	// avoid any kind of leaks
	defer cursor.Close(ctx)

	var notes []Note

	if err := cursor.All(ctx, &notes); err != nil {
		return nil, fmt.Errorf("Decode notes failed: %w", err)
	}

	return notes, nil
}

func (r *Repo) GetByID(ctx context.Context, id primitive.ObjectID) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	filter := bson.M{"_id": id}

	var note Note

	err := r.coll.FindOne(opCtx, filter, options.FindOne()).Decode(&note)

	if err != nil {
		return Note{}, fmt.Errorf("Find note by id failed: %w", err)
	}
	return note, nil
}

func (r *Repo) UpdateByID(ctx context.Context, id primitive.ObjectID, req UpdateNoteRequest) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filer := bson.M{"_id": id}

	update := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"pinned":    req.Pinned,
			"UpdatedAt": time.Now().UTC(),
		},
	}

	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var updated Note

	err := r.coll.FindOneAndUpdate(opCtx, filer, update, &opts).Decode(&updated)

	if err != nil {
		return Note{}, fmt.Errorf("Update note failed: %w", err)
	}

	return updated, nil

}
