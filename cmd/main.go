package main

import (
	"context"
	"errors"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.NewConfig()

	sessionTTL := 7 * 24 * time.Hour
	accessTTL := 5 * time.Minute

	runMigrations(cfg.DatabaseURL)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

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

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.JSONLogger(logger))

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
		slog.Info("server started", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()

	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited properly")
}

func runMigrations(databaseURL string) {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		slog.Error("failed to create migrate instance", "error", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations applied")
}
