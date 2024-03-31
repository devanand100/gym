package store

import (
	"context"
	"errors"
	"time"

	"github.com/devanand100/gym/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Gym struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AddressId primitive.ObjectID `bson:"AddressId,omitempty"`
	CompanyId primitive.ObjectID `bson:"companyId"`

	GymName   string `bson:"gymName"`
	ContactNo int64  `bson:"contactNo"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Deleted   bool      `bson:"deleted"`
}

func (s Store) CreateGym(ctx context.Context, gym *dto.GymCreateReq) (primitive.ObjectID, error) {
	var err error

	var addressId primitive.ObjectID
	addressId, err = s.CreateAddress(ctx, Address{City: gym.City, PinCode: gym.PinCode, StreetAddress: gym.StreetAddress})
	if err != nil {
		return primitive.NilObjectID, err
	}

	gymCollection := s.DB.Collection("gym")

	newGym := &Gym{
		GymName:   gym.GymName,
		CompanyId: gym.CompanyObjectId,
		AddressId: addressId,
		ContactNo: gym.ContactNo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}

	var createdGym *mongo.InsertOneResult
	createdGym, err = gymCollection.InsertOne(ctx, newGym)

	if err != nil {
		return primitive.NilObjectID, err
	}

	InsertedID, ok := createdGym.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("New Company InsertedID is not ObjectID type")
	}

	return InsertedID, nil
}

func (s Store) UpdateGym(ctx context.Context, gym *dto.GymUpdateReq) error {
	var err error
	gymCollection := s.DB.Collection("gym")

	err = s.UpdateAddress(ctx, Address{City: gym.City, PinCode: gym.PinCode, StreetAddress: gym.StreetAddress, ID: gym.AddressObjectId})

	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": gym.GymObjectId,
	}

	update := bson.M{
		"$set": bson.M{
			"gymName":   gym.GymName,
			"contactNo": gym.ContactNo,
		},
	}

	_, err = gymCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil

}

func (s Store) GetGymById(ctx context.Context, gymId primitive.ObjectID) (gymDoc *Gym, err error) {
	gymCollection := s.DB.Collection("gym")

	err = gymCollection.FindOne(ctx, bson.M{"_id": gymId}).Decode(&gymDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}

	return
}

func (s Store) GetGymList(ctx context.Context) (gyms []*Gym, err error) {

	gymCollection := s.DB.Collection("gym")

	var cursor *mongo.Cursor
	cursor, err = gymCollection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var gym Gym
		err = cursor.Decode(&gym)
		if err != nil {
			return nil, err
		}
		gyms = append(gyms, &gym)
	}

	return
}
