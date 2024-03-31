package api

import (
	"encoding/json"
	"net/http"

	"github.com/devanand100/gym/internal/dto"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *APIService) RegisterGymRoutes(g *echo.Group) {
	g.GET("/gym", s.GetGymList)
	g.POST("/gym", s.CreateGym)
	g.PUT("/gym/:gymId", s.UpdateGym)
	g.GET("/gym/:gymId", s.GetGymById)
}

func (s *APIService) CreateGym(c echo.Context) error {
	ctx := c.Request().Context()
	var err error

	gymCreate := &dto.GymCreateReq{}

	if err := json.NewDecoder(c.Request().Body).Decode(gymCreate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed post user request").SetInternal(err)
	}
	if err := gymCreate.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid gym create format").SetInternal(err)
	}

	var CompanyObjectId primitive.ObjectID
	CompanyObjectId, err = primitive.ObjectIDFromHex(gymCreate.CompanyId)
	gymCreate.CompanyObjectId = CompanyObjectId

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Company Id").SetInternal(err)
	}

	var gymId primitive.ObjectID
	gymId, err = s.Store.CreateGym(ctx, gymCreate)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gymId)
}

func (s *APIService) UpdateGym(c echo.Context) (err error) {
	ctx := c.Request().Context()

	gymUpdate := &dto.GymUpdateReq{}

	if err := json.NewDecoder(c.Request().Body).Decode(gymUpdate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Malformed PUT user request").SetInternal(err)
	}
	if err := gymUpdate.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid gym update format").SetInternal(err)
	}

	gymId := c.Param("gymId")

	var GymObjectId primitive.ObjectID
	GymObjectId, err = primitive.ObjectIDFromHex(gymId)
	gymUpdate.GymObjectId = GymObjectId

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Gym Id").SetInternal(err)
	}

	var gym *store.Gym
	gym, err = s.Store.GetGymById(ctx, GymObjectId)

	if err != nil {
		return err
	}

	if gym == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Gym Not Found").SetInternal(err)
	}

	gymUpdate.AddressObjectId = gym.AddressId
	err = s.Store.UpdateGym(ctx, gymUpdate)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Gym Updated Successfully")
}

func (s *APIService) GetGymById(c echo.Context) (err error) {

	ctx := c.Request().Context()
	gymId := c.Param("gymId")

	var gymObjectId primitive.ObjectID
	gymObjectId, err = primitive.ObjectIDFromHex(gymId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Gym Id").SetInternal(err)
	}

	var gym *store.Gym
	gym, err = s.Store.GetGymById(ctx, gymObjectId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gym)
}

func (s *APIService) GetGymList(c echo.Context) (err error) {
	ctx := c.Request().Context()

	gyms, err := s.Store.GetGymList(ctx)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gyms)
}
