package main

import (
	"log"
	"os"
	"planner/config"
	"planner/internal/auth"
	"planner/internal/handler"
	"planner/internal/service"
	"planner/internal/store"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
	authService := service.NewAuthService(s, a.Tokens)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	public := r.Group("/")
	protected := r.Group("/api", auth.Middleware(a.Tokens))

	authHandler.RegisterRoutes(public, protected)

	log.Println("Server starting on :" + strconv.Itoa(cfg.Server.Port))
	r.Run(":" + strconv.Itoa(cfg.Server.Port))
}
