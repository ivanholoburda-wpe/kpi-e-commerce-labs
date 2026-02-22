package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"

	"lab2/internal/api/handlers"
	"lab2/internal/api/middleware"
	"lab2/internal/auth"
	"lab2/internal/clients/telegram"
	"lab2/internal/config"
	"lab2/internal/db/sqlc"
	"lab2/internal/repository"
	"lab2/internal/services"
)

func main() {
	cfg := config.NewConfig()

	sessionTTL := 7 * 24 * time.Hour
	accessTTL := 5 * time.Minute

	runMigrations(cfg.DatabaseURL)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %s", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping database: %s", err)
	}
	log.Println("Connected to database")

	queries := sqlc.New(pool)

	accountRepo := repository.NewAccountRepo(queries)
	codeRepo := repository.NewVerificationCodeRepo(queries)
	tgClient := telegram.NewClient(cfg.TelegramToken)

	registerService := services.NewRegisterService(accountRepo, codeRepo, tgClient, cfg.JWTSecret)
	authService := services.NewAuthService(accountRepo, cfg.JWTSecret, sessionTTL, accessTTL)
	recoveryService := services.NewRecoveryService(accountRepo, codeRepo, tgClient, cfg.JWTSecret)

	registerHandler := handlers.NewRegisterHandler(registerService)
	authHandler := handlers.NewAuthHandler(authService)
	recoveryHandler := handlers.NewRecoveryHandler(recoveryService)

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.AbortWithError(http.StatusServiceUnavailable, err)
			return
		}

		c.Status(http.StatusOK)
	})

	reg := router.Group("/register")
	{
		reg.POST("/start", registerHandler.Start)
		reg.POST("/verify", registerHandler.Verify)

		regAuth := reg.Group("", middleware.RequireJWT(cfg.JWTSecret, auth.PurposeRegistration))
		{
			regAuth.POST("/setup", registerHandler.Setup)
			regAuth.POST("/confirm-totp", registerHandler.ConfirmTOTP)
		}
	}

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)

		authGroup.POST("/verify-2fa",
			middleware.RequireJWT(cfg.JWTSecret, auth.Purpose2FA),
			authHandler.Verify2FA,
		)
		authGroup.POST("/setup-pin",
			middleware.RequireJWT(cfg.JWTSecret, auth.PurposePinSetup),
			authHandler.SetupPIN,
		)
		authGroup.POST("/verify-pin",
			middleware.RequireJWT(cfg.JWTSecret, auth.PurposeSession),
			authHandler.VerifyPIN,
		)
		authGroup.GET("/me",
			middleware.RequireJWT(cfg.JWTSecret, auth.PurposeAccess),
			authHandler.GetMe,
		)
	}

	rec := router.Group("/recovery")
	{
		rec.POST("/start", recoveryHandler.Start)
		rec.POST("/verify-otp", recoveryHandler.VerifyOTP)

		rec.POST("/confirm-pin",
			middleware.RequireJWT(cfg.JWTSecret, auth.PurposeRecoveryOTP),
			recoveryHandler.ConfirmPIN,
		)
		rec.POST("/reset-password",
			middleware.RequireJWT(cfg.JWTSecret, auth.PurposeRecoveryConfirmed),
			recoveryHandler.ResetPassword,
		)
	}

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s...\n", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen error: %s\n", err)
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()

	log.Println("SIGTERM received. Starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exited properly")
}

func runMigrations(databaseURL string) {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		log.Fatalf("create migrate instance: %s", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("run migrations: %s", err)
	}

	log.Println("Migrations applied successfully")
}
