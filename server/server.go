package server

import (
	"context"
	"fmt"
	"time"

	api "github.com/devanand100/gym/api"
	"github.com/devanand100/gym/server/profile"
	"github.com/devanand100/gym/store"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo

	Profile *profile.Profile
	Store   *store.Store
}

func NewServer(ctx context.Context, profile *profile.Profile, store *store.Store) (*Server, error) {

	e := echo.New()
	e.Debug = true

	s := &Server{
		e:       e,
		Store:   store,
		Profile: profile,
	}

	apiService := api.NewApiService(profile, store)

	rootGroup := e.Group("")

	apiService.Register(rootGroup)
	return s, nil
}

func (s *Server) Start(ctx context.Context) error {

	fmt.Println("server start")
	return s.e.Start(fmt.Sprintf("%s:%d", s.Profile.Addr, s.Profile.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.e.Shutdown(ctx); err != nil {
		fmt.Println("Server ShutDown Error")
	}
	return nil
}
