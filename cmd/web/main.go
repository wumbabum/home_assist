package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/wumbabum/home_assist/internal/authenticator"
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
	auth0 struct {
		domain       string
		clientID     string
		clientSecret string
		callbackURL  string
	}
	baseURL  string
	httpPort int
	db       struct {
		dsn         string
		automigrate bool
	}
	session struct {
		cookieName string
	}
}

type application struct {
	auth0          *authenticator.Authenticator
	config         config
	db             *database.DB
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
}

func run(logger *slog.Logger) error {
	// Register types for session storage
	gob.Register(UserProfile{})

	// Configure Environment
	var cfg config

	cfg.auth0.domain = env.GetString("AUTH0_DOMAIN", "placeholder-domain.auth0.com")
	cfg.auth0.clientID = env.GetString("AUTH0_CLIENT_ID", "placeholder-client-id")
	cfg.auth0.clientSecret = env.GetString("AUTH0_CLIENT_SECRET", "placeholder-client-secret")
	cfg.auth0.callbackURL = env.GetString("AUTH0_CALLBACK_URL", "http://localhost:5749/callback")
	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:5749")
	cfg.httpPort = env.GetInt("HTTP_PORT", 5749)
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
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}()

	if cfg.db.automigrate {
		err = db.MigrateUp()
		if err != nil {
			return err
		}
	}

	auth0, err := authenticator.New()
	if err != nil {
		return err
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db.DB.DB)
	sessionManager.Lifetime = 7 * 24 * time.Hour
	sessionManager.Cookie.Name = cfg.session.cookieName
	sessionManager.Cookie.Secure = true

	app := &application{
		auth0:          auth0,
		config:         cfg,
		db:             db,
		logger:         logger,
		sessionManager: sessionManager,
	}

	return app.serveHTTP()
}
