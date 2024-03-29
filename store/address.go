package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	City          string             `bson:"city"`
	StreetAddress string             `bson:"streetAddress"`
	PinCode       int64              `bson:"password"`
	Country       string             `bson:"country"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Deleted   bool      `bson:"deleted"`
}

func (s Store) CreateAddress(ctx context.Context, address Address) (primitive.ObjectID, error) {
	addressCollection := s.DB.Collection("address")

	newAddress := &Address{
		PinCode:       address.PinCode,
		City:          address.City,
		Country:       address.Country,
		StreetAddress: address.StreetAddress,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}

	add, err := addressCollection.InsertOne(ctx, newAddress)
	if err != nil {
		return primitive.NilObjectID, err
	}

	InsertedID, ok := add.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("New Address InsertedID is not ObjectID type")
	}

	return InsertedID, nil
}
