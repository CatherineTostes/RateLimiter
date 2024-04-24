package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	middleware2 "github.com/catherinetostes/ratelimiter/internal/middleware"
	"github.com/catherinetostes/ratelimiter/internal/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestPerIP(t *testing.T) {
	os.Chdir("../..")
	t.Run("execute only one request per ip", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		rateLimiter := middleware2.NewRateLimiter()
		r.Use(rateLimiter.RateLimiterPerClient)
		r.Get("/exchange", web.HandleExchange)
		server := httptest.NewServer(r)
		defer server.Close()

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/exchange", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("execute several requests per ip that will be blocked by the rate limit", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		rateLimiter := middleware2.NewRateLimiter()
		r.Use(rateLimiter.RateLimiterPerClient)
		r.Get("/exchange", web.HandleExchange)
		server := httptest.NewServer(r)
		defer server.Close()

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/exchange", nil)

		for i := 0; i < 11; i++ {
			resp, _ := client.Do(req)
			defer resp.Body.Close()

			if i > 10 {
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}

	})
}

func TestGetRequestPerAPITOKEN(t *testing.T) {
	t.Run("execute only one request per API_TOKEN", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		rateLimiter := middleware2.NewRateLimiter()
		r.Use(rateLimiter.RateLimiterPerClient)
		r.Get("/exchange", web.HandleExchange)
		server := httptest.NewServer(r)
		defer server.Close()

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/exchange", nil)
		req.Header.Set("API_TOKEN", "abc123")
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("execute several requests per ip that will be blocked by the rate limit", func(t *testing.T) {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		rateLimiter := middleware2.NewRateLimiter()
		r.Use(rateLimiter.RateLimiterPerClient)
		r.Get("/exchange", web.HandleExchange)
		server := httptest.NewServer(r)
		defer server.Close()

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/exchange", nil)
		req.Header.Set("API_TOKEN", "abc123")

		for i := 0; i < 13; i++ {
			resp, _ := client.Do(req)
			defer resp.Body.Close()

			if i > 12 {
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}

	})
}
