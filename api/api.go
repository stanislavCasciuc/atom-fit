package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/stanislavCasciuc/atom-fit/api/handlers"
	customMiddleware "github.com/stanislavCasciuc/atom-fit/api/middleware"
	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/auth"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type userKey string

const UserCtx userKey = "user"

type Application struct {
	Config config.Config
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
	authenticator := auth.New(a.Config.Auth.Secret, a.Config.Auth.Aud)
	h := handlers.New(resp, a.Store, a.Config, authenticator)
	m := customMiddleware.New(a.Store, resp, authenticator)
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
				r.Put("/activate", h.ActivateUser)
				r.Post("/register", h.RegisterUserHandler)
				r.Post("/login", h.LoginHandler)
				r.With(m.AuthTokenMiddleware).Get("/", h.GetUserHandler)
				r.Route("/attributes", func(r chi.Router) {
					r.Use(m.AuthTokenMiddleware)
					r.Get("/", h.GetUserWithAttrHandler)
					r.Post("/log/weight/", h.LogWeightHandler)
				})
			})
		})
	})
	return r
}
