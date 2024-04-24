package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/catherinetostes/ratelimiter/cmd/configs"
	"github.com/catherinetostes/ratelimiter/db"
	"github.com/redis/go-redis/v9"
)

var message = Message{
	Body: "You have reached the maximum number of requests or actions allowed within a certain time frame",
}

type (
	RateLimiter struct {
		redisConnection *db.RedisConnection
	}

	Message struct {
		Status string `json:"status"`
		Body   string `json:"body"`
	}
)

func NewRateLimiter() *RateLimiter {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	database, err := strconv.Atoi(configs.RedisDatabase)
	if err != nil {
		panic(err)
	}

	redisConnection, err := db.NewRedisConnection(configs.RedisAddr, configs.RedisPassword, database)
	if err != nil {
		panic(err)
	}
	return &RateLimiter{redisConnection: redisConnection}
}

func (rl *RateLimiter) RateLimiterPerClient(next http.Handler) http.Handler {
	var counter int64

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token := r.Header.Get("API_TOKEN")

		var key string
		if token != "" {
			key = fmt.Sprintf("API_TOKEN:%s", token)
		} else {
			key = fmt.Sprintf("IP:%s", ip)
		}

		var timeDuration int

		if strings.Contains(key, "API_TOKEN") {
			timeDuration, err = strconv.Atoi(configs.BlockingTimeToken)
			if err != nil {
				panic(err)
			}

			counter, err = rl.redisConnection.Get(r.Context(), key)

			if err == redis.Nil {
				err = rl.redisConnection.Set(r.Context(), key, "1", time.Duration(timeDuration)*time.Second)
				if err != nil {
					panic(err)
				}
				counter = 1
			}

			limitRequestPerToken, err := strconv.ParseInt(configs.LimitRequestPerToken, 10, 64)
			if err != nil {
				panic(err)
			}

			if counter > limitRequestPerToken {
				http.Error(w, message.Body, http.StatusTooManyRequests)
				return
			}
		} else {
			timeDuration, err = strconv.Atoi(configs.BlockingTimeIP)
			if err != nil {
				panic(err)
			}

			counter, err = rl.redisConnection.Get(r.Context(), key)

			if err == redis.Nil {
				err = rl.redisConnection.Set(r.Context(), key, "1", time.Duration(timeDuration)*time.Second)
				if err != nil {
					panic(err)
				}
				counter = 1
			}

			limitRequestPerIP, err := strconv.ParseInt(configs.LimitRequestPerIP, 10, 64)
			if err != nil {
				panic(err)
			}

			if counter > limitRequestPerIP {
				http.Error(w, message.Body, http.StatusTooManyRequests)
				return
			}
		}

		counter, err = rl.redisConnection.Incr(r.Context(), key)

		next.ServeHTTP(w, r)
	})
}
