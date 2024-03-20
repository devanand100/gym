package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devanand100/gym/api/auth"
	"github.com/devanand100/gym/internal/util"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIService) RegisterUserRoutes(g *echo.Group) {
	g.POST("/user/register", s.RegisterUser)
	g.POST("/user/login", s.LoginUser)
	g.GET("/user/me", s.Me)
}

func (s *APIService) RegisterUser(c echo.Context) error {
	// TODO : send verification mail
	ctx := c.Request().Context()
	userRegister := &RegisterReq{}
	if err := json.NewDecoder(c.Request().Body).Decode(userRegister); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed post user request").SetInternal(err)
	}
	if err := userRegister.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user create format").SetInternal(err)
	}

	user, err := s.Store.FindUserByEmail(ctx, userRegister.Email)

	if err != nil && err != mongo.ErrNoDocuments {
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

	return c.JSON(http.StatusOK, "User Registered successfully")
}

func (s *APIService) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()
	userLogin := &LoginReq{}

	if err := json.NewDecoder(c.Request().Body).Decode(userLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed post user request").SetInternal(err)
	}

	user, err := s.Store.FindUserByEmail(ctx, userLogin.Email)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect Email or Password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		fmt.Println("err............", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect Email or Password").SetInternal(err)
	}

	token, err := auth.CreateToken(user.ID.String())

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Token Generation Error")
	}

	cookie := http.Cookie{
		Name:     auth.CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(&cookie)
	return c.String(http.StatusOK, "Login successful")
}

func (s *APIService) Me(c echo.Context) error {
	cookie, err := c.Cookie(auth.CookieName)

	if err != nil {
		if err == http.ErrNoCookie {
			return echo.NewHTTPError(http.StatusUnauthorized, "Cookie not found")
		}
		return err
	}
	token := cookie.Value

	if _, err := auth.VerifyToken(token); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token Expired").SetInternal(err)
	}

	return c.String(http.StatusOK, "Auth Success")
}

type RegisterReq struct {
	FirstName string `json:"firstName" `
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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
