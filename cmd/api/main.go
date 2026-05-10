package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/masamodelkin/nudge-server/config"
	"github.com/masamodelkin/nudge-server/internal/auth"
	"github.com/masamodelkin/nudge-server/internal/handler"
	"github.com/masamodelkin/nudge-server/internal/service"
	"github.com/masamodelkin/nudge-server/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := sqlx.Connect("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	db.MustExec("PRAGMA foreign_keys = ON;")

	s := store.New(db)
	a := auth.New(
		os.Getenv("JWT_SECRET"),
		cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
	)

	r := gin.Default()
	public := r.Group("/")
	protected := r.Group("/api", auth.Middleware(a.Tokens))

	authService := service.NewAuthService(s, a.Tokens)
	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(public, protected)

	statusService := service.NewStatusService(s)
	statusHandler := handler.NewStatusHandler(statusService)
	statusHandler.RegisterRoutes(protected)

	labelsService := service.NewLabelService(s)
	labelsHandler := handler.NewLabelHandler(labelsService)
	labelsHandler.RegisterRoutes(protected)

	log.Println("Server starting on :" + strconv.Itoa(cfg.Server.Port))
	r.Run(":" + strconv.Itoa(cfg.Server.Port))
}
