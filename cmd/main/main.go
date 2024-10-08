package main

import (
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/stanislavCasciuc/atom-fit/api"
	"github.com/stanislavCasciuc/atom-fit/db"
	"github.com/stanislavCasciuc/atom-fit/internal/env"
	"github.com/stanislavCasciuc/atom-fit/internal/lib/config"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

//	@title			Atom Fit API
//	@description	This is a sample server for Atom Fit API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/api/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @tokenUrl					/auth/login
// @description
func main() {
	iatEnvString := env.EnvString("IAT", "24h")
	iatDuratioin, _ := time.ParseDuration(iatEnvString)

	cfg := config.Config{
		Addr:         env.EnvString("ADDR", ":8080"),
		CompleteAddr: env.EnvString("COMPL_ADDR", "https://localhost:8080"),
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
		Auth: config.Auth{
			Secret: env.EnvString("SECRET", "secret"),
			Aud:    env.EnvString("AUD", "atom-fit"),
			Iat:    iatDuratioin,
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
