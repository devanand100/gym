package dto

import (
	"github.com/devanand100/gym/internal/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyCreateReq struct {
	CompanyName   string `json:"companyName" `
	City          string `json:"city"`
	StreetAddress string `json:"streetAddress"`
	PinCode       int64  `json:"pinCode"`
	Country       string `json:"country"`
	//Company status Live or Test
	Status string `json:"status"`
	// country code of mobile no
	ContactNo      int64  `json:"contactNo"`
	OwnerFirstName string `json:"ownerFirstName"`
	OwnerLastName  string `json:"ownerLastName"`
	OwnerEmail     string `json:"ownerEmail"`
}

type CompanyUpdateReq struct {
	CompanyId         string `json:"companyId"`
	CompanyIdObjectId primitive.ObjectID
	AddressId         primitive.ObjectID
	CompanyName       string `json:"companyName" `
	City              string `json:"city"`
	StreetAddress     string `json:"streetAddress"`
	PinCode           int64  `json:"pinCode"`
	Country           string `json:"country"`
	//Company status Live or Test
	Status string `json:"status"`
	// country code of mobile no
	ContactNo int64 `json:"contactNo"`
}

func (company CompanyCreateReq) Validate() error {
	if len(company.CompanyName) > 100 {
		return errors.New("Company Name To long")
	}

	if company.OwnerEmail != "" {
		if len(company.OwnerEmail) > 256 {
			return errors.New("email is too long, maximum length is 256")
		}

		if !util.ValidateEmail(company.OwnerEmail) {
			return errors.New("invalid email format")
		}
	}

	if company.City == "" {
		return errors.New("City is required")
	}

	if company.StreetAddress == "" {
		return errors.New("Street Address is required")
	}

	if company.PinCode == 0 {
		return errors.New("Pin Code is required")
	}

	if company.Status == "" {
		return errors.New("Company Status Is Required")
	}
	if company.Status != "LIVE" && company.Status != "TEST" {
		return errors.New("Invalid Company Status ")
	}

	if company.Country == "" {
		return errors.New("Country is required")
	}

	if !isValidMobileNumber(company.ContactNo) {
		return errors.New("Contact No is invalid")
	}

	if company.OwnerFirstName == "" {
		return errors.New("Owner First Name is required")
	}

	if company.OwnerLastName == "" {
		return errors.New("Owner Last Name is required")
	}

	return nil
}

func (company CompanyUpdateReq) Validate() error {
	if len(company.CompanyName) > 100 {
		return errors.New("Company Name To long")
	}

	if company.City == "" {
		return errors.New("City is required")
	}

	if company.StreetAddress == "" {
		return errors.New("Street Address is required")
	}

	if company.PinCode == 0 {
		return errors.New("Pin Code is required")
	}

	if company.Status != "LIVE" && company.Status != "TEST" {
		return errors.New("Invalid Company Status ")
	}

	if company.Country == "" {
		return errors.New("Country is required")
	}

	if !isValidMobileNumber(company.ContactNo) {
		return errors.New("Contact No is invalid")
	}

	return nil
}

func isValidMobileNumber(number int64) bool {

	return number > 0 && number >= 1000000000 && number <= 9999999999
}
