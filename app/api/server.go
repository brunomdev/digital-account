package api

import (
	"errors"
	"github.com/brunomdev/digital-account/app/api/middleware"
	appConfig "github.com/brunomdev/digital-account/config"
	"github.com/brunomdev/digital-account/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"net/http"
)

type Server struct {
	cfg        *appConfig.Config
	newRelic   *newrelic.Application
	httpServer *fiber.App
	service    *domain.Service
}

func NewServer(options ...func(server *Server) error) (*Server, error) {
	server := &Server{}
	for _, option := range options {
		err := option(server)
		if err != nil {
			return nil, err
		}
	}

	app := fiber.New()

	middleware.FiberMiddleware(app, server.newRelic, server.cfg.AppDebug)

	server.httpServer = app

	server.router()

	go func() {
		if err := server.httpServer.Listen(":" + server.cfg.HTTPPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return server, nil
}

func WithConfig(cfg *appConfig.Config) func(server *Server) error {
	return func(server *Server) error {
		server.cfg = cfg
		return nil
	}
}

func WithService(service *domain.Service) func(server *Server) error {
	return func(server *Server) error {
		server.service = service
		return nil
	}
}

func WithNewRelic(nr *newrelic.Application) func(server *Server) error {
	return func(server *Server) error {
		server.newRelic = nr
		return nil
	}
}

func (s *Server) Close() error {
	return s.httpServer.Shutdown()
}
