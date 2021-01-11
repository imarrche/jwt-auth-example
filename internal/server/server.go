package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"

	"github.com/imarrche/jwt-auth-example/internal/config"
	"github.com/imarrche/jwt-auth-example/internal/logger"
	"github.com/imarrche/jwt-auth-example/internal/service"
)

// Server is REST API server.
type Server struct {
	server  *http.Server
	router  *chi.Mux
	service service.Service
}

// New returns new Server instance.
func New(c *config.Server, service service.Service) *Server {
	r := chi.NewRouter()
	s := &http.Server{
		Addr:         c.Addr,
		Handler:      r,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	return &Server{server: s, router: r, service: service}
}

// Run runs the server.
func (s *Server) Run() {
	// Middleware setup.
	s.router.Use(s.loggerMiddleware())

	// Router setup.
	s.configureRouter()

	// Graceful shutdown setup.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	// Starting the server.
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Get().Fatal("couldn't start the server")
		}
	}()
	logger.Get().Info("server started")

	<-done
	// Gracefully shutting down the server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Get().Fatal("couldn't shut down server gracefully")
	}

	logger.Get().Info("server shutted down gracefully")
}

// configureRouter maps all handlers.
func (s *Server) configureRouter() {
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", s.signUp())
			r.Post("/sign-in", s.signIn())
			r.Post("/refresh", s.refresh())
		})

		r.Get("/public", s.public())
		r.Route("/private", func(r chi.Router) {
			r.Use(s.authMiddleware())
			r.Get("/", s.private())
		})
	})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if err != nil {
		s.respond(w, r, code, map[string]string{"error": err.Error()})
		return
	}
	s.respond(w, r, code, nil)
}
