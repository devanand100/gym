package dto

import (
	"github.com/pkg/errors"

	"github.com/devanand100/gym/internal/util"
)

type UserType string

const (
	GlobalUserType  UserType = "GLOBAL"
	CompanyUserType UserType = "COMPANY"
	GymUserType     UserType = "GYM"
)

type RegisterReq struct {
	FirstName   string `json:"firstName" `
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	HasPassword string
	UserType    UserType
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (register RegisterReq) Validate() error {
	if len(register.Password) < 3 {
		return errors.New("password is too short, minimum length is 3")
	}
	if len(register.Password) > 512 {
		return errors.New("password is too long, maximum length is 512")
	}
	if len(register.FirstName) > 64 {
		return errors.New("firstName is too long, maximum length is 64")
	}
	if len(register.LastName) > 64 {
		return errors.New("lastName is too long, maximum length is 64")
	}
	if register.Email != "" {
		if len(register.Email) > 256 {
			return errors.New("email is too long, maximum length is 256")
		}
		if !util.ValidateEmail(register.Email) {
			return errors.New("invalid email format")
		}
	}

	return nil
}

func (login LoginReq) Validate() error {

	if len(login.Email) == 0 {
		return errors.New("Password is Required")
	}

	if !util.ValidateEmail(login.Email) {
		return errors.New("invalid email format")
	}

	if len(login.Password) == 0 {
		return errors.New("Password is Required")
	}
	return nil
}
