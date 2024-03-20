package api

import (
	"github.com/devanand100/gym/api/auth"
	"github.com/devanand100/gym/server/profile"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
)

type APIService struct {
	Profile *profile.Profile
	Store   *store.Store
}

func NewApiService(profile *profile.Profile, store *store.Store) *APIService {
	return &APIService{
		Profile: profile,
		Store:   store,
	}
}

func (s *APIService) Register(rootGroup *echo.Group) {
	apiGroup := rootGroup.Group("/api")

	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return JWTMiddleware(s, next, string(auth.SecretKey))
	})
	s.RegisterSystemRoutes(apiGroup)
	s.RegisterUserRoutes(apiGroup)
	s.RegisterCompanyRoutes(apiGroup)
}
