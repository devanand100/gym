package server

import (
	"context"
	"fmt"

	"github.com/devanand100/gym/api"
	"github.com/devanand100/gym/server/profile"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	e *echo.Echo

	// Secret  string
	Profile *profile.Profile
	Client  *mongo.Client
}

func NewServer(ctx context.Context, profile *profile.Profile, client *mongo.Client) (*Server, error) {

	e := echo.New()
	e.Debug = true

	s := &Server{
		e:       e,
		Client:  client,
		Profile: profile,
	}

	apiService := api.NewApiService(profile, client)

	rootGroup := e.Group("")

	apiService.Register(rootGroup)
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {

	fmt.Println("server start")
	return s.e.Start(fmt.Sprintf("%s:%d", s.Profile.Addr, s.Profile.Port))
}
