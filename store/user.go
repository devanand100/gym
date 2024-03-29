package store

import (
	"context"
	"errors"
	"time"

	"github.com/devanand100/gym/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FirstName   string             `bson:"firstName"`
	LastName    string             `bson:"lastName"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	UserType    dto.UserType       `bson:"userType"`
	Permissions []string           `bson:"permissions"`
	RoleId      primitive.ObjectID `bson:"roleId,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	Deleted     bool               `bson:"deleted"`
}

func (s *Store) RegisterUser(ctx context.Context, user *dto.RegisterReq) (primitive.ObjectID, error) {
	userCollection := s.DB.Collection("user")

	newUser := User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.HasPassword,
		CreatedAt:   time.Now(),
		Permissions: []string{},
		RoleId:      primitive.NilObjectID,
		UpdatedAt:   time.Now(),
		Deleted:     false,
		UserType:    user.UserType,
	}

	registeredUser, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return primitive.NilObjectID, err
	}

	InsertedID, ok := registeredUser.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("New Address InsertedID is not ObjectID type")
	}

	return InsertedID, nil

}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	userCollection := s.DB.Collection("user")
	filter := bson.M{"email": email}
	var user User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
