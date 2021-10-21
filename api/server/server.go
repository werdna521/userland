package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/handler/auth"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
	"github.com/werdna521/userland/service"
)

type Server struct {
	Config
	DataSource   *DataSource
	repositories *repositories
	services     *services
}

type repositories struct {
	ur repository.UserRepository
}

type services struct {
	as service.AuthService
}

type Config struct {
	Port string
}

type DataSource struct {
	Postgres *sql.DB
}

func NewServer(config Config, dataSource *DataSource) *Server {
	return &Server{
		Config:     config,
		DataSource: dataSource,
	}
}

func (s *Server) Start() {
	log.Info().Msg("initializing repositories")
	s.initRepositories()
	defer s.tearDownRepositories()

	log.Info().Msg("initializing services")
	s.initServices()

	log.Info().Msg("initializing handlers")
	h := s.initHandlers()
	port := fmt.Sprintf(":%s", s.Port)

	log.Info().Msgf("server running on port %s", port)
	http.ListenAndServe(port, h)
}

func (s *Server) initRepositories() {
	ur := postgres.NewUserRepository(s.DataSource.Postgres)
	ur.PrepareStatements(context.Background())

	s.repositories = &repositories{
		ur: ur,
	}
}

func (s *Server) tearDownRepositories() {
	s.repositories.ur.TearDownStatements()
}

func (s *Server) initServices() {
	as := service.NewBaseAuthService(s.repositories.ur)

	s.services = &services{
		as: as,
	}
}

func (s *Server) initHandlers() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(s.services.as))
		})
	})

	return r
}
