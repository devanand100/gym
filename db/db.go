package db

import (
	"context"

	"github.com/devanand100/gym/server/profile"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client  *mongo.Client
	Profile *profile.Profile
	GymDB   *mongo.Database
}

func NewDb(ctx context.Context, profile *profile.Profile) (*DB, error) {

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(profile.DbUri))

	GymDB := client.Database("gym")

	if err != nil {
		return nil, err
	}

	return &DB{
		client:  client,
		Profile: profile,
		GymDB:   GymDB,
	}, nil
}

// CloseDB closes the database connection and cancels the context.
func (db *DB) Close(ctx context.Context) error {

	if err := db.client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
