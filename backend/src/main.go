package src

import (
	"context"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	wallet "wallet/src/modules/wallet"
	wallet_model "wallet/src/modules/wallet/model"
)

type Env struct {
	HTTP_Host          string `env:"HTTP_HOST" default:"localhost"`
	HTTP_Port          int    `env:"HTTP_PORT" default:"80"`
	POSTGRES_Host      string `env:"POSTGRES_HOST"`
	POSTGRES_User      string `env:"POSTGRES_USER"`
	POSTGRES_Password  string `env:"POSTGRES_PASSWORD"`
	POSTGRES_DBName    string `env:"POSTGRES_DB"`
	POSTGRES_Port      int    `env:"POSTGRES_PORT"`
	POSTGRES_SSLMode   string `env:"POSTGRES_SSL_MODE" default:"disable"`
	POSTGRES_TimeZone  string `env:"POSTGRES_TIME_ZONE" default:"UTC"`
	JWT_Secret         string `env:"JWT_SECRET"`
	JWT_Expired        int64  `env:"JWT_EXPIRED"`
	JWT_RefreshExpired int64  `env:"JWT_REFRESH_EXPIRED"`
	REDIS_Addr         string `env:"REDIS_ADDR"`
	TON_Network        string `env:"TON_NETWORK" default:"testnet"`
	ENCRYPTION_KEY     string `env:"ENCRYPTION_KEY"`
	DEBUG              bool   `env:"DEBUG" default:"false"`
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// @title           TON Wallet API
// @version         1.0
// @description     REST API for managing TON blockchain wallets
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey Bearer
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and JWT token.
func Exec(env *Env) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		env.POSTGRES_Host,
		env.POSTGRES_User,
		env.POSTGRES_Password,
		env.POSTGRES_DBName,
		env.POSTGRES_Port,
		env.POSTGRES_SSLMode,
		env.POSTGRES_TimeZone,
	)

	gormLogLevel := logger.Silent
	if env.DEBUG {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migration(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = &customValidator{validator: validator.New()}

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	e.Static("/swagger", "./static/swagger")
	e.File("/swagger.json", "./docs/swagger/swagger.json")
	e.File("/swagger.yaml", "./docs/swagger/swagger.yaml")

	if err := wallet.Cmd(e, db, env.TON_Network, env.ENCRYPTION_KEY); err != nil {
		log.Fatalf("Failed to initialize wallet module: %v", err)
	}

	if err := e.Start(fmt.Sprintf(":%d", env.HTTP_Port)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func migration(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&wallet_model.User{},
		&wallet_model.Wallet{},
		&wallet_model.Transaction{},
	); err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	return db.WithContext(context.Background()).Exec(`
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE OR REPLACE FUNCTION HASH_MAKE(password TEXT) RETURNS TEXT AS $$
		BEGIN
			RETURN crypt(password, gen_salt('bf'));
		END;
		$$ LANGUAGE plpgsql;

		CREATE OR REPLACE FUNCTION HASH_CHECK(password TEXT, hashed_password TEXT) RETURNS BOOLEAN AS $$
		BEGIN
			RETURN crypt(password, hashed_password) = hashed_password;
		END;
		$$ LANGUAGE plpgsql;
	`).Error
}
