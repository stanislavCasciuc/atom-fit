package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/stanislavCasciuc/atom-fit/api/handlers"
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type Config struct {
	Addr string
	DB   DbConfig
	Env  string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
type Application struct {
	Config Config
	Log    *zap.SugaredLogger
	Store  store.Storage
}

func (a *Application) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         a.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	a.Log.Infow("server has started", "addr", a.Config.Addr, "env", a.Config.Env)

	return srv.ListenAndServe()
}

func (a *Application) Mount() http.Handler {
	resp := response.New(a.Log)
	h := handlers.New(resp, a.Store)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/health", h.HealthHandler)
			r.Route("/users", func(r chi.Router) {
				r.Post("/register", h.RegisterUserHandler)
			})
		})
	})
	return r
}
