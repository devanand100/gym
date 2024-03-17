package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIService) RegisterUserRoutes(g *echo.Group) {
	g.POST("/user/register", s.RegisterUser)
}

func (*APIService) RegisterUser(c echo.Context) error {

	return c.JSON(http.StatusOK, "pong")
}
