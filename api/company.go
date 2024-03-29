package api

import (
	"encoding/json"
	"net/http"

	"github.com/devanand100/gym/dto"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *APIService) RegisterCompanyRoutes(g *echo.Group) {
	g.POST("/company", s.CreateCompany)
}

func (s *APIService) CreateCompany(c echo.Context) error {
	ctx := c.Request().Context()
	var err error

	companyCreate := &dto.CompanyCreateReq{}

	if err := json.NewDecoder(c.Request().Body).Decode(companyCreate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed post user request").SetInternal(err)
	}
	if err := companyCreate.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid company create format").SetInternal(err)
	}

	existingUser, err := s.Store.FindUserByEmail(ctx, companyCreate.OwnerEmail)

	if err != nil {
		return err
	}

	if existingUser != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "This Email is Already registered").SetInternal(err)
	}

	var companyId primitive.ObjectID
	companyId, err = s.Store.CreateCompany(ctx, companyCreate)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, companyId)
}
