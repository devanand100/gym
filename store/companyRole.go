package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyRole struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CompanyId primitive.ObjectID `bson:"companyId" `

	Permissions []string `bson:"permissions"`
	RoleName    string   `bson:"roleName"`
	Description string   `bson:"description"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Deleted   bool      `bson:"deleted"`
}

func (s Store) CreateCompanyRole(ctx context.Context, role *CompanyRoleCreate) (primitive.ObjectID, error) {

	companyRoleCollection := s.DB.Collection("companyRole")

	newRole := CompanyRole{
		CompanyId:   role.CompanyId,
		RoleName:    role.RoleName,
		Description: role.Description,
		Permissions: role.Permissions,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}

	createdRole, err := companyRoleCollection.InsertOne(ctx, newRole)
	if err != nil {
		return primitive.NilObjectID, err
	}

	InsertedID, ok := createdRole.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("New Address InsertedID is not ObjectID type")
	}

	return InsertedID, nil
}

type CompanyRoleCreate struct {
	CompanyId   primitive.ObjectID
	RoleName    string
	Description string
	Permissions []string
}
