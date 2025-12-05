package src

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	wallet "wallet_test/src/modules/wallet"
	walletModel "wallet_test/src/modules/wallet/model"
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
	TON_Network        string `env:"TON_NETWORK" default:"testnet"` // mainnet или testnet
	ENCRYPTION_KEY     string `env:"ENCRYPTION_KEY"`                // Ключ для шифрования seed фраз (32 байта)
}

// @title TON Wallet API
// @version 1.0
// @description REST API для управления TON блокчейн кошельками
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func Exec(env *Env) {
	// Initialize Bun database connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		env.POSTGRES_User,
		env.POSTGRES_Password,
		env.POSTGRES_Host,
		env.POSTGRES_Port,
		env.POSTGRES_DBName,
		env.POSTGRES_SSLMode,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging (optional)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	defer db.Close()

	// Run migrations
	if err := migration(db); err != nil {
		log.Println("Migration error:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	setupRoutes(router, db, env)

	// Swagger - отдаем статические файлы из docs/swagger
	router.Static("/swagger", "./docs/swagger")

	// Start server
	router.Run(fmt.Sprintf(":%d", env.HTTP_Port))
}

func migration(db *bun.DB) error {
	ctx := context.Background()

	// Create tables
	models := []interface{}{
		(*walletModel.User)(nil),
		(*walletModel.Wallet)(nil),
		(*walletModel.Transaction)(nil),
	}

	for _, m := range models {
		_, err := db.NewCreateTable().Model(m).IfNotExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Create PostgreSQL functions
	_, err := db.ExecContext(ctx, `
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
	`)

	return err
}

func setupRoutes(router *gin.Engine, db *bun.DB, env *Env) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Wallet routes
	wallet.Cmd(router, db, env.TON_Network, env.ENCRYPTION_KEY)
}
