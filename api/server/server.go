package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/handler/auth"
	"github.com/werdna521/userland/api/handler/session"
	"github.com/werdna521/userland/api/handler/user"
	"github.com/werdna521/userland/api/middleware"
	"github.com/werdna521/userland/mailer"
	"github.com/werdna521/userland/repository/postgres"
	rds "github.com/werdna521/userland/repository/redis"
	"github.com/werdna521/userland/service"
)

type Server struct {
	Config
	mailer       mailer.Mailer
	DataSource   *DataSource
	repositories *repositories
	services     *services
}

type repositories struct {
	ur  postgres.UserRepository
	phr postgres.PasswordHistoryRepository
	tr  rds.TokenRepository
	sr  rds.SessionRepository
}

type services struct {
	as service.AuthService
	ss service.SessionService
	us service.UserService
}

type Config struct {
	Port string
}

type DataSource struct {
	Postgres *sql.DB
	Redis    *redis.Client
}

func NewServer(config Config, mailer mailer.Mailer, dataSource *DataSource) *Server {
	return &Server{
		Config:     config,
		mailer:     mailer,
		DataSource: dataSource,
	}
}

func (s *Server) Start() {
	log.Info().Msg("initializing repositories")
	s.initRepositories()

	log.Info().Msg("initializing services")
	s.initServices()

	log.Info().Msg("initializing handlers")
	h := s.initHandlers()
	port := fmt.Sprintf(":%s", s.Port)

	log.Info().Msgf("server running on port %s", port)
	http.ListenAndServe(port, h)
}

func (s *Server) initRepositories() {
	ur := postgres.NewBaseUserRepository(s.DataSource.Postgres)
	ur.PrepareStatements(context.Background())

	phr := postgres.NewBasePasswordHistoryRepository(s.DataSource.Postgres)
	phr.PrepareStatements(context.Background())

	tr := rds.NewBaseTokenRepository(s.DataSource.Redis)

	sr := rds.NewBaseSessionRepository(s.DataSource.Redis)

	s.repositories = &repositories{
		ur:  ur,
		phr: phr,
		tr:  tr,
		sr:  sr,
	}
}

func (s *Server) initServices() {
	as := service.NewBaseAuthService(
		s.repositories.ur,
		s.repositories.phr,
		s.repositories.tr,
		s.repositories.sr,
		s.mailer,
	)

	ss := service.NewBaseSessionService(s.repositories.sr)

	us := service.NewBaseUserService(
		s.repositories.ur,
		s.repositories.phr,
		s.repositories.tr,
		s.repositories.sr,
		s.mailer,
	)

	s.services = &services{
		as: as,
		ss: ss,
		us: us,
	}
}

func (s *Server) initHandlers() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(s.services.as))
			r.Post("/login", auth.Login(s.services.as))

			r.Route("/verification", func(r chi.Router) {
				r.Get("/", auth.VerifyEmail(s.services.as))
				r.Post("/", auth.SendVerification(s.services.as))
			})

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", auth.ForgotPassword(s.services.as))
				r.Post("/reset", auth.ResetPassword(s.services.as))
			})
		})

		r.Route("/me", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Get("/", user.GetInfoDetail(s.services.us))
				r.Post("/", user.UpdateBasicInfo(s.services.us))
			})

			r.Route("/email", func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Get("/", user.GetCurrentEmailAddress(s.services.us))
				r.Post("/", user.RequestEmailAddressChange(s.services.us))
			})

			r.Route("/password", func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Post("/", user.ChangePassword(s.services.us))
			})

			r.Route("/picture", func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Post("/", user.SetProfilePicture(s.services.us))
				r.Delete("/", user.DeleteProfilePicture(s.services.us))
			})

			r.Route("/delete", func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Post("/", user.DeleteAccount(s.services.us))
			})

			r.Route("/session", func(r chi.Router) {
				r.Use(middleware.ValidateAccessToken(s.repositories.sr))

				r.Get("/", session.ListSessions(s.services.ss))
				r.Delete("/", session.EndCurrentSession(s.services.ss))
				r.Delete("/other", session.DeleteAllOtherSessions(s.services.ss))
				r.Post("/refresh_token", session.GenerateRefreshToken(s.services.ss))
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.ValidateRefreshToken(s.repositories.sr))
				r.Post("/session/access_token", session.GenerateAccessToken(s.services.ss))
			})

			r.Group(func(r chi.Router) {
				r.Get("/email/verification", user.VerifyEmailChange(s.services.us))
			})
		})
	})

	s.initFileServer(r)

	return r
}

func (s *Server) initFileServer(r chi.Router) {
	workDir, _ := os.Getwd()
	// get the path to the static files directory
	path := filepath.Join(workDir, "uploaded")

	// create the directory if it doesn't exist
	os.MkdirAll(path, os.ModePerm)

	// create a file server for the static files
	fileServer(r, "/uploaded", http.Dir(path))
}
