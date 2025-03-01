package main

import (
	"github.com/dankru/API_gateway/internal/middleware"
	"github.com/dankru/API_gateway/internal/proxy"
	"log"
	"net/http"
	"time"
)

func main() {
	timeout := time.Second * 10
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			MaxConnsPerHost:    100,
			DisableKeepAlives:  false,
		},
	}

	proxyHandler := proxy.NewProxy(timeout, client)
	rateLimiter := middleware.NewSlidingWindowLimiter(time.Minute, 100)

	proxyHandler.AddRoute("/api/users/", "http://localhost:8080/users/")
	proxyHandler.AddRoute("/auth/", "http://localhost:8080/auth/")

	handler := rateLimiter.Middleware(proxyHandler)
	server := &http.Server{
		Addr:    ":3000",
		Handler: handler,
	}

	log.Println("starting server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}
