package store

import (
	"context"
	"errors"
	"time"

	"github.com/devanand100/gym/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Company struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AddressId primitive.ObjectID `bson:"AddressId,omitempty"`

	CompanyName string `bson:"companyName" `
	ContactNo   int64  `bson:"contactNo"`

	//owner Details
	OwnerFirstName string `bson:"ownerFirstName"`
	OwnerLastName  string `bson:"ownerLastName"`
	OwnerEmail     string `bson:"ownerEmail"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Deleted   bool      `bson:"deleted"`
}

func (s Store) CreateCompany(ctx context.Context, company *dto.CompanyCreateReq) (primitive.ObjectID, error) {
	var err error

	var addressId primitive.ObjectID
	addressId, err = s.CreateAddress(ctx, Address{City: company.City, PinCode: company.PinCode, StreetAddress: company.StreetAddress, Country: company.Country})
	if err != nil {
		return primitive.NilObjectID, err
	}

	companyCollection := s.DB.Collection("company")

	newCompany := &Company{
		CompanyName:    company.CompanyName,
		AddressId:      addressId,
		ContactNo:      company.ContactNo,
		OwnerFirstName: company.OwnerFirstName,
		OwnerLastName:  company.OwnerLastName,
		OwnerEmail:     company.OwnerEmail,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Deleted:        false,
	}

	var createdCompany *mongo.InsertOneResult
	createdCompany, err = companyCollection.InsertOne(ctx, newCompany)

	if err != nil {
		return primitive.NilObjectID, err
	}

	InsertedID, ok := createdCompany.InsertedID.(primitive.ObjectID)

	if !ok {
		return primitive.NilObjectID, errors.New("New Company InsertedID is not ObjectID type")
	}

	//Create default Roles for new Company
	var adminRoleId primitive.ObjectID
	adminRoleId, err = s.CreateCompanyRole(ctx, &CompanyRoleCreate{CompanyId: InsertedID, RoleName: "Admin", Description: "Admin User that is main Company User and that have all company level changes permission ", Permissions: []string{"admin"}})
	if err != nil {
		return primitive.NilObjectID, err
	}
	_, err = s.CreateCompanyRole(ctx, &CompanyRoleCreate{CompanyId: InsertedID, RoleName: "Manager", Description: "Company User that manages company", Permissions: []string{"manager"}})
	if err != nil {
		return primitive.NilObjectID, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost) //TODO:Change this constant password to random generated password
	if err != nil {
		return primitive.NilObjectID, err
	}

	//registering company Owner as User and Admin User of Company
	var userId primitive.ObjectID
	userId, err = s.RegisterUser(ctx, &dto.RegisterReq{
		FirstName:   company.OwnerFirstName,
		LastName:    company.OwnerLastName,
		Email:       company.OwnerEmail,
		HasPassword: string(hashPassword),
		UserType:    dto.CompanyUserType,
	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	err = s.CreateCompanyUser(ctx, &CompanyUserCreate{CompanyId: InsertedID, RoleId: adminRoleId, UserId: userId})
	// TODO:send Invite mail and with new Password creation
	if err != nil {
		return primitive.NilObjectID, err
	}

	return InsertedID, nil
}

type CompanyCreate struct {
	CompanyName   string `json:"companyName" `
	City          string `json:"city"`
	StreetAddress string `json:"streetAddress"`
	PinCode       int64  `json:"pinCode"`
	Country       string `json:"country"`
	// country code of mobile no
	ContactNo      int64  `json:"contactNo"`
	OwnerFirstName string `json:"ownerFirstName"`
	OwnerLastName  string `json:"ownerLastName"`
	OwnerEmail     string `json:"ownerEmail"`
}
