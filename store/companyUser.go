package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CompanyId primitive.ObjectID `bson:"companyId" `
	RoleId    primitive.ObjectID `bson:"roleId" `
	//reference of user table
	UserId primitive.ObjectID `bson:"userId" `

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Deleted   bool      `bson:"deleted"`
}

func (s Store) CreateCompanyUser(ctx context.Context, user *CompanyUserCreate) error {

	companyUserCollection := s.DB.Collection("companyUser")

	newUser := CompanyUser{
		CompanyId: user.CompanyId,
		RoleId:    user.RoleId,
		UserId:    user.UserId,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}

	_, err := companyUserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

type CompanyUserCreate struct {
	CompanyId primitive.ObjectID
	RoleId    primitive.ObjectID
	UserId    primitive.ObjectID
}
