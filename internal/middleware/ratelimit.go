package middleware

import (
	"net/http"
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	windowSize time.Duration
	limit      int
	requests   map[string][]time.Time
	mu         sync.RWMutex
}

func NewSlidingWindowLimiter(windowSize time.Duration, limit int) *SlidingWindowLimiter {
	limiter := &SlidingWindowLimiter{
		windowSize: windowSize,
		limit:      limit,
		requests:   make(map[string][]time.Time),
	}

	return limiter
}

func (l *SlidingWindowLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
