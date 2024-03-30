package api

import (
	"encoding/json"
	"net/http"

	"github.com/devanand100/gym/internal/dto"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *APIService) RegisterCompanyRoutes(g *echo.Group) {
	g.POST("/company", s.CreateCompany)
	g.PUT("/company/:companyId", s.UpdateCompany)
	g.GET("/company/:companyId", s.GetCompanyById)
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

func (s *APIService) UpdateCompany(c echo.Context) (err error) {
	ctx := c.Request().Context()

	companyUpdate := &dto.CompanyUpdateReq{}

	if err := json.NewDecoder(c.Request().Body).Decode(companyUpdate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed PUT user request").SetInternal(err)
	}
	if err := companyUpdate.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid company update format").SetInternal(err)
	}

	companyId := c.Param("companyId")

	var companyIdObjectId primitive.ObjectID
	companyIdObjectId, err = primitive.ObjectIDFromHex(companyId)
	companyUpdate.CompanyIdObjectId = companyIdObjectId
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Company Id").SetInternal(err)
	}
	var company *store.Company
	company, err = s.Store.GetCompanyById(ctx, companyIdObjectId)

	if err != nil {
		return err
	}

	if company == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Company Not Found").SetInternal(err)
	}

	companyUpdate.AddressId = company.AddressId
	err = s.Store.UpdateCompany(ctx, companyUpdate)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Company Updated Successfully")
}

func (s *APIService) GetCompanyById(c echo.Context) (err error) {

	ctx := c.Request().Context()
	companyId := c.Param("companyId")

	var companyIdObjectId primitive.ObjectID
	companyIdObjectId, err = primitive.ObjectIDFromHex(companyId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Company Id").SetInternal(err)
	}

	var company *store.Company
	company, err = s.Store.GetCompanyById(ctx, companyIdObjectId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, company)
}
