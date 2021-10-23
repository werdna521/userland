package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/handler/auth"
	"github.com/werdna521/userland/repository"
	"github.com/werdna521/userland/repository/postgres"
	rds "github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/service"
)

type Server struct {
	Config
	DataSource   *DataSource
	repositories *repositories
	services     *services
}

type repositories struct {
	ur  repository.UserRepository
	evr repository.EmailVerificationRepository
	fpr repository.ForgotPasswordRepository
}

type services struct {
	as service.AuthService
}

type Config struct {
	Port string
}

type DataSource struct {
	Postgres *sql.DB
	Redis    *redis.Client
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

	evr := rds.NewVerificationRepository(s.DataSource.Redis)

	fpr := rds.NewForgotPasswordRepository(s.DataSource.Redis)

	s.repositories = &repositories{
		ur:  ur,
		evr: evr,
		fpr: fpr,
	}
}

func (s *Server) tearDownRepositories() {
	s.repositories.ur.TearDownStatements()
}

func (s *Server) initServices() {
	as := service.NewBaseAuthService(
		s.repositories.ur,
		s.repositories.evr,
		s.repositories.fpr,
	)

	s.services = &services{
		as: as,
	}
}

func (s *Server) initHandlers() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(s.services.as))

			r.Route("/verification", func(r chi.Router) {
				r.Get("/", auth.VerifyEmail(s.services.as))
				r.Post("/", auth.SendVerification(s.services.as))
			})

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", auth.ForgotPassword(s.services.as))
			})
		})
	})

	return r
}
