package api

import (
	"github.com/devanand100/gym/server/profile"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIService struct {
	Profile *profile.Profile
	Client  *mongo.Client
}

func NewApiService(profile *profile.Profile, client *mongo.Client) *APIService {
	return &APIService{
		Profile: profile,
		Client:  client,
	}
}

func (s *APIService) Register(rootGroup *echo.Group) {
	apiGroup := rootGroup.Group("/api")

	s.RegisterSystemRoutes(apiGroup)
}
