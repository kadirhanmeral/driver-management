package handlers

import (
	"fmt"
	"go-api-gateway/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func MakeProxyHandler(serviceName string, targetBase string, cfg *config.Config, cache *cache.Cache) gin.HandlerFunc {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		log.Fatalf("Invalid URL for %s: %v", serviceName, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if isRateLimited(serviceName, clientIP, cfg.RateLimitWindow, cfg.RateLimitCount, cache) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		originalPath := c.Request.URL.Path
		newPath := originalPath

		c.Request.URL.Path = newPath
		c.Request.URL.RawPath = newPath
		c.Request.Host = targetURL.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func isRateLimited(serviceName, ip string, window time.Duration, limit int, cache *cache.Cache) bool {
	key := fmt.Sprintf("%s-%s", serviceName, ip)

	val, err := cache.IncrementInt(key, 1)
	if err != nil {
		cache.Set(key, 1, window)
		return false
	}

	if val > limit {
		return true
	}

	return false
}
