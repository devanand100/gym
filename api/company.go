package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *APIService) RegisterCompanyRoutes(g *echo.Group) {
	g.POST("/company", s.CreateCompany)
}

func (s *APIService) CreateCompany(c echo.Context) error {
	return c.JSON(http.StatusOK, "Company Created ")
}
