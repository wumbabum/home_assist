package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/wumbabum/home_assist/internal/database"
	"github.com/wumbabum/home_assist/internal/env"
	"github.com/wumbabum/home_assist/internal/version"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	session struct {
		cookieName string
	}
}

type application struct {
	config         config
	db             *database.DB
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:5749")
	cfg.httpPort = env.GetInt("HTTP_PORT", 5749)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")
	cfg.db.dsn = env.GetString("DB_DSN", "user:pass@localhost:5432/db")
	cfg.db.automigrate = env.GetBool("DB_AUTOMIGRATE", true)
	cfg.session.cookieName = env.GetString("SESSION_COOKIE_NAME", "session_ux762yqp")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if cfg.db.automigrate {
		err = db.MigrateUp()
		if err != nil {
			return err
		}
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db.DB.DB)
	sessionManager.Lifetime = 7 * 24 * time.Hour
	sessionManager.Cookie.Name = cfg.session.cookieName
	sessionManager.Cookie.Secure = true

	app := &application{
		config:         cfg,
		db:             db,
		logger:         logger,
		sessionManager: sessionManager,
	}

	return app.serveHTTP()
}
