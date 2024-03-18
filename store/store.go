package store

import (
	"github.com/devanand100/gym/server/profile"
	"go.mongodb.org/mongo-driver/mongo"
)

// Store provides database access to all
type Store struct {
	Profile *profile.Profile
	DB      *mongo.Database
}

// New creates a new instance of Store.
func New(db *mongo.Database, profile *profile.Profile) *Store {

	return &Store{
		DB:      db,
		Profile: profile,
	}
}
