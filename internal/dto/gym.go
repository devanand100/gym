package dto

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GymCreateReq struct {
	CompanyId       string `json:"companyId"`
	GymName         string `json:"gymName"`
	ContactNo       int64  `json:"contactNo"`
	City            string `json:"city"`
	StreetAddress   string `json:"streetAddress"`
	PinCode         int64  `json:"pinCode"`
	CompanyObjectId primitive.ObjectID
}

type GymUpdateReq struct {
	GymName       string `json:"gymName"`
	ContactNo     int64  `json:"contactNo"`
	City          string `json:"city"`
	StreetAddress string `json:"streetAddress"`
	PinCode       int64  `json:"pinCode"`
	//used for updating address
	AddressObjectId primitive.ObjectID
	GymObjectId     primitive.ObjectID
}

func (gym GymCreateReq) Validate() error {
	if gym.CompanyId == "" {
		return errors.New("Company Id is required")
	}

	if len(gym.GymName) > 100 {
		return errors.New("Gym Name To long")
	}

	if gym.City == "" {
		return errors.New("City is required")
	}

	if gym.StreetAddress == "" {
		return errors.New("Street Address is required")
	}

	if gym.PinCode == 0 {
		return errors.New("Pin Code is required")
	}

	if !isValidMobileNumber(gym.ContactNo) {
		return errors.New("Contact No is invalid")
	}

	return nil
}

func (gym GymUpdateReq) Validate() error {
	if len(gym.GymName) > 100 {
		return errors.New("Gym Name To long")
	}

	if gym.City == "" {
		return errors.New("City is required")
	}

	if gym.StreetAddress == "" {
		return errors.New("Street Address is required")
	}

	if gym.PinCode == 0 {
		return errors.New("Pin Code is required")
	}

	if !isValidMobileNumber(gym.ContactNo) {
		return errors.New("Contact No is invalid")
	}

	return nil
}
