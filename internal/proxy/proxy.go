package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	routes  map[string]*url.URL
	client  *http.Client
	timeout time.Duration
	rewrite map[string]string
}

func NewProxy(timeout time.Duration, client *http.Client) *Proxy {
	return &Proxy{
		routes:  make(map[string]*url.URL),
		client:  client,
		timeout: timeout,
		rewrite: make(map[string]string),
	}
}

func (p *Proxy) AddRoute(prefix string, backend string) error {
	backendURL, err := url.Parse(backend)
	if err != nil {
		return err
	}
	p.routes[prefix] = backendURL
	return nil
}
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var targetURL *url.URL
	var longestPrefix string

	for prefix, backend := range p.routes {
		if len(prefix) > len(longestPrefix) && strings.HasPrefix(r.URL.Path, prefix) {
			longestPrefix = prefix
			targetURL = backend
		}
	}

	if targetURL == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}

	outReq := p.CreateProxyRequest(r, targetURL)
	resp, err := p.client.Do(outReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute request: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf(
		"Proxy request: %s %s -> %s, status: %d, latency: %v",
		r.Method,
		r.URL.Path,
		targetURL.String(),
		resp.StatusCode,
		time.Since(startTime),
	)
}

func (p *Proxy) CreateProxyRequest(r *http.Request, target *url.URL) *http.Request {
	outURL := *target // Начинаем с базового backend URL

	idx := strings.Index(r.URL.Path, target.Path)
	if idx != -1 {
		outURL.Path = r.URL.Path[idx:]
	}

	// Объединяем query параметры
	if r.URL.RawQuery != "" && target.RawQuery != "" {
		outURL.RawQuery = r.URL.RawQuery + "&" + target.RawQuery
	} else {
		outURL.RawQuery = r.URL.RawQuery + target.RawQuery
	}

	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, outURL.String(), r.Body)
	if err != nil {
		log.Printf("Error creating proxy request: %s", err.Error())
		return nil
	}

	// Копируем заголовки
	for key, values := range r.Header {
		for _, value := range values {
			outReq.Header.Add(key, value)
		}
	}

	outReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	outReq.Header.Set("X-Forwarded-Host", r.Host)
	outReq.Header.Set("X-Forwarded-Proto", r.URL.Scheme)

	return outReq
}
