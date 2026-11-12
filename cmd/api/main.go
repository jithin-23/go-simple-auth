package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/config"
	"simple-auth/internal/database"
	"simple-auth/internal/server"
	"simple-auth/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	setupLogger(cfg.GinMode)
	slog.Info("config loaded",
		"port", cfg.Port, "gin_mode", cfg.GinMode,
		"db", cfg.DB.Host+":"+cfg.DB.Port+"/"+cfg.DB.Name, "jwt_ttl", cfg.JWT.TTL.String())

	db, err := database.Connect(cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("database connected")

	userRepo := user.NewGormRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	gin.SetMode(cfg.GinMode)
	router := server.NewRouter(db, server.Handlers{User: userHandler})

	slog.Info("server starting", "addr", ":"+cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

// setupLogger uses readable text logs locally and JSON logs in production.
func setupLogger(ginMode string) {
	var handler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if ginMode == gin.ReleaseMode {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(handler))
}
