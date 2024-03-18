package api

import (
	"encoding/json"
	"net/http"

	"github.com/devanand100/gym/internal/util"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIService) RegisterUserRoutes(g *echo.Group) {
	g.POST("/user/register", s.RegisterUser)
}

func (s *APIService) RegisterUser(c echo.Context) error {
	// TODO : send verification mail
	ctx := c.Request().Context()
	userRegister := &RegisterReq{}
	if err := json.NewDecoder(c.Request().Body).Decode(userRegister); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformatted post user request").SetInternal(err)
	}
	if err := userRegister.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user create format").SetInternal(err)
	}

	user, err := s.Store.FindUserByEmail(ctx, userRegister.Email)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Find user by email error").SetInternal(err)
	}

	if user != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "This email is already registered")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password hash").SetInternal(err)
	}

	err = s.Store.RegisterUser(ctx, &store.RegisterUser{
		FirstName:   userRegister.FirstName,
		LastName:    userRegister.LastName,
		Email:       userRegister.Email,
		HasPassword: string(hashPassword),
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user").SetInternal(err)
	}

	return c.JSON(http.StatusOK, nil)
}

type RegisterReq struct {
	FirstName string `json:"firstName" `
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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
