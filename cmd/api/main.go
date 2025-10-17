package main

import (
	"time"

	"github.com/MohammadTaghipour/social/internal/auth"
	"github.com/MohammadTaghipour/social/internal/db"
	"github.com/MohammadTaghipour/social/internal/env"
	"github.com/MohammadTaghipour/social/internal/mailer"
	"github.com/MohammadTaghipour/social/internal/store"
	"go.uber.org/zap"
)

const version string = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for Gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		}, mail: mailConfig{
			mailHog: mailHogConfig{
				addr: env.GetString("MAILHOG_ADDR", "localhost:1025"),
			},
			fromEmail: env.GetString("FROM_EMAIL", "gopher@email.com"),
			exp:       time.Hour * 24 * 3, // 3 days to accept invitations
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			jwt: jwtConfig{
				secret: env.GetString("AUTH_JWT_SECRET", "supersecretkey"),
				issuer: env.GetString("AUTH_JWT_ISSUER", "gophersocial"),
				expiration: time.Duration(
					env.GetInt("AUTH_JWT_EXPIRATION_HOURS", 72)) * time.Hour,
			},
		},
		env: env.GetString("ENV", "dev"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
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

	store := store.NewStorage(db) // TODO: pass a real db connection

	mailer := mailer.NewMailhog(cfg.mail.mailHog.addr, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJwtAuthenticator(
		cfg.auth.jwt.secret,
		cfg.auth.jwt.issuer,
		cfg.auth.jwt.issuer,
	)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
