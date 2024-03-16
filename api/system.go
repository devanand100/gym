package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIService) RegisterSystemRoutes(g *echo.Group) {
	g.GET("/ping", s.PingSystem)
}

func (*APIService) PingSystem(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
