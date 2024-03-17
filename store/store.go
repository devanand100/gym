package store

import (
	"github.com/devanand100/gym/server/profile"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store provides database access to all
type Store struct {
	Profile *profile.Profile
	db      *mongo.Client
}

// New creates a new instance of Store.
func New(db *mongo.Client, profile *profile.Profile) *Store {
	return &Store{
		db:      db,
		Profile: profile,
	}
}
