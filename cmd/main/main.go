package main

import (
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/stanislavCasciuc/atom-fit/api"
	"github.com/stanislavCasciuc/atom-fit/db"
	"github.com/stanislavCasciuc/atom-fit/internal/env"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

func main() {
	cfg := config.Config{
		Addr: env.EnvString("ADDR", ":8080"),
		DB: config.DbConfig{
			Addr: env.EnvString(
				"DB_ADDR",
				"postgres://postgres:postgres@localhost:5432/atom-fit?sslmode=disable",
			),
			MaxOpenConns: env.IntEnv("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.IntEnv("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.EnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.EnvString("ENV", "dev"),
		Mail: config.MailCfg{
			Addr:     env.EnvString("EMAIL_ADDR", ""),
			Host:     env.EnvString("EMAIL_HOST", ""),
			Port:     env.IntEnv("EMAIL_PORT", 0),
			Password: env.EnvString("EMAIL_PASS", ""),
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("db connected successfully")

	store := store.New(db)

	app := &api.Application{
		Config: cfg,
		Log:    logger,
		Store:  store,
	}

	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
