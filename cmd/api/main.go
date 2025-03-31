package main

import (
	"expvar"
	"fmt"
	"runtime"

	"github.com/salimofshadow/usenet-client/internal/routes/db"
	"github.com/salimofshadow/usenet-client/internal/routes/env"
	"github.com/salimofshadow/usenet-client/internal/routes/store"

	"go.uber.org/zap"
)

const version = "0.0.1"

type application struct {
	config        config
	store         store.Storage
	// cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	// mailer        mailer.Client
	// authenticator auth.Authenticator
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	// mail        mailConfig
	frontendURL string
	// auth        authConfig
	// redisCfg    redisConfig
	// rateLimiter ratelimiter.Config
}

func main() {
	fmt.Println("cfg")
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/postgres"),			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Main Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	

	store := store.NewStorage(db)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
	}

	// Metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}