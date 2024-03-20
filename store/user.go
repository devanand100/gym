package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUser struct {
	FirstName   string
	LastName    string
	Email       string
	HasPassword string
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FirstName   string             `bson:"firstName"`
	LastName    string             `bson:"lastName"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Permissions []string           `bson:"permissions"`
	RoleId      primitive.ObjectID `bson:"roleId,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	Deleted     bool               `bson:"deleted"`
}

func (s *Store) RegisterUser(ctx context.Context, user *RegisterUser) error {
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
	}

	_, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return err
	}

	return nil

}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	userCollection := s.DB.Collection("user")
	filter := bson.M{"email": email}
	var user User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
