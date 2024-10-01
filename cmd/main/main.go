package main

import (
	"go.uber.org/zap"

	"github.com/stanislavCasciuc/atom-fit/api"
	"github.com/stanislavCasciuc/atom-fit/internal/env"
)

func main() {
	config := api.Config{
		Addr: env.EnvString("ADDR", ":8080"),
		DB: api.DbConfig{
			Addr: env.EnvString(
				"DB_ADDR",
				"postgres://postgres:postgres@localhost:5432/atom-fit?sslmode=disable",
			),
			MaxOpenConns: env.IntEnv("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.IntEnv("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.EnvString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.EnvString("ENV", "dev"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	app := &api.Application{
		Config: config,
		Log:    logger,
	}

	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
