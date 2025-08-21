package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	// "time"

	"github.com/Android-Shubham/auth/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	// "github.com/robfig/cron/v3"
)

type ApiConfig struct {
	db *database.Queries
	secret string
}

func main() {

	// c := cron.New()
	// c.AddFunc("@every 1m", func() {
	// 	// Add your cron job logic here
	// 	fmt.Println("Cron job executed at ", time.Now())
	// })
	// c.Start()


	fmt.Println("Hello, World!")
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Error loading .env file")
		return
	}
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		fmt.Println("Error loading DB_URL from .env file")
		return
	}

	conn, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	queries := database.New(conn)
	apiConfig := ApiConfig{
		db: queries,
		secret: "let_it_be_for_n0w", // This should be replaced with a more secure secret in production
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum age in seconds for preflight requests
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", apiConfig.checkHealthHandler)
	v1Router.Get("/error", apiConfig.errorHandler)
	v1Router.Post("/users", apiConfig.createUser)
	v1Router.Post("/login", apiConfig.loginUser)
	v1Router.Get("/user", apiConfig.middlewareAuth(apiConfig.getUserDetails))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
