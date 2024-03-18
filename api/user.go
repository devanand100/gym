package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIService) RegisterUserRoutes(g *echo.Group) {
	g.POST("/user/register", s.RegisterUser)
}

func (s *APIService) RegisterUser(c echo.Context) error {
	// TODO : send verification mail
	return c.JSON(http.StatusOK, "pong")
}
