package main

import (
	"fmt"
	"go-api-gateway/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

const (
	ENV          = "dev"
	SERVICE_PORT = "8080"
)

func main() {
	router := mux.NewRouter()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	cache := cache.New(5*time.Minute, cfg.RateLimitWindow)

	registerServices(router, cfg, cache)

	log.Printf("API Gateway running on :%s", SERVICE_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVICE_PORT, router))
}

func registerServices(router *mux.Router, cfg *config.Config, cache *cache.Cache) {

	for serviceName, svc := range cfg.Services {
		for _, route := range svc.Routes {
			// /driver-service/drivers -> /drivers
			// We register the full path in the gateway: /serviceName + route.Path
			gatewayPath := "/" + serviceName + route.Path

			log.Printf("Registering route: %s -> %s%s", gatewayPath, svc.BaseURL, route.Path)

			router.Handle(gatewayPath, makeProxyHandler(serviceName, svc.BaseURL, route.Path, cfg, cache))
		}
	}
}

func makeProxyHandler(serviceName string, targetBase string, targetPath string, cfg *config.Config, cache *cache.Cache) http.Handler {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		log.Fatalf("Invalid URL for %s: %v", serviceName, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		clientIP := r.RemoteAddr

		if isRateLimited(serviceName, clientIP, cfg.RateLimitWindow, cfg.RateLimitCount, cache) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// gorilla/mux parametreleri al
		vars := mux.Vars(r)

		// backend'e gidecek yeni path
		newPath := targetPath

		// targetPath = "/drivers/{id}" gibi ise -> {id} yerini değiştir
		for k, v := range vars {
			placeholder := "{" + k + "}"
			newPath = strings.ReplaceAll(newPath, placeholder, v)
		}

		// örn: /drivers/{id} -> /drivers/692ec3bfde4b612ec30a6647
		r.URL.Path = newPath
		r.URL.RawPath = newPath

		log.Printf("[%s] %s -> %s%s", serviceName, r.RequestURI, targetURL, newPath)

		proxy.ServeHTTP(w, r)
	})

}

func isRateLimited(serviceName, ip string, window time.Duration, limit int, cache *cache.Cache) bool {
	key := fmt.Sprintf("%s-%s", serviceName, ip)

	// Increment returns an error if the item doesn't exist
	val, err := cache.IncrementInt(key, 1)
	if err != nil {
		// Item doesn't exist, set it with initial value 1 and expiration
		cache.Set(key, 1, window)
		return false
	}

	if val > limit {
		return true
	}

	return false
}
