package main

import (
	"fmt"
	"net/http"

	"github.com/catherinetostes/ratelimiter/cmd/configs"
	middleware2 "github.com/catherinetostes/ratelimiter/internal/middleware"
	"github.com/catherinetostes/ratelimiter/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting web server on port", configs.WebServerPort)

	rateLimiter := middleware2.NewRateLimiter()

	route := chi.NewRouter()
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(rateLimiter.RateLimiterPerClient)
	route.Get("/exchange", web.HandleExchange)

	http.ListenAndServe(":8080", route)

	fmt.Println("Server running on port 8080")
}
