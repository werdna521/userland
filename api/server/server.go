package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/handler/auth"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
)

type Server struct {
	Config
	DataSource   *DataSource
	repositories *repositories
}

type repositories struct {
	ur repository.UserRepository
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

	log.Info().Msg("initializing handlers")
	h := s.initHandlers()
	port := fmt.Sprintf(":%s", s.Port)

	log.Info().Msgf("server running on port %s", port)
	http.ListenAndServe(port, h)
}

func (s *Server) initRepositories() {
	ur := postgres.NewUserRepository(s.DataSource.Postgres)
	s.repositories = &repositories{
		ur: ur,
	}
}

func (s *Server) initHandlers() http.Handler {
	r := chi.NewRouter()

	r.Get("/", auth.HandleRegister(s.repositories.ur))

	return r
}
