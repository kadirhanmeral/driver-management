package handlers

import (
	"encoding/json"
	"go-api-gateway/config"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Param("any")

		switch path {
		case "/doc.json":
			gatewayDoc, err := os.ReadFile("docs/swagger.json")
			if err != nil {
				log.Printf("Failed to load gateway swagger: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load gateway swagger"})
				return
			}

			var merged map[string]interface{}
			if err := json.Unmarshal(gatewayDoc, &merged); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse gateway swagger"})
				return
			}

			driverSvc, ok := cfg.Services["driver"]
			if ok {
				targetURL := driverSvc.BaseURL + "/swagger/doc.json"
				log.Printf("Fetching Swagger from: %s", targetURL)
				resp, err := http.Get(targetURL)
				if err != nil {
					log.Printf("Failed to fetch driver swagger: %v", err)
				} else if resp.StatusCode != 200 {
					log.Printf("Failed to fetch driver swagger: status %d", resp.StatusCode)
					resp.Body.Close()
				} else {
					defer resp.Body.Close()
					var driverDoc map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&driverDoc); err != nil {
						log.Printf("Failed to decode driver swagger: %v", err)
					} else {
						log.Printf("Successfully fetched driver swagger. Merging...")
						if paths, ok := driverDoc["paths"].(map[string]interface{}); ok {
							if mergedPaths, ok := merged["paths"].(map[string]interface{}); ok {
								for pathKey, pathValue := range paths {
									if pathItem, ok := pathValue.(map[string]interface{}); ok {
										for method, methodValue := range pathItem {
											if methodOp, ok := methodValue.(map[string]interface{}); ok {
												methodOp["security"] = []map[string][]string{
													{"BearerAuth": {}},
												}
												pathItem[method] = methodOp
											}
										}
										mergedPaths[pathKey] = pathItem
									}
								}
							} else {
								for pathKey, pathValue := range paths {
									if pathItem, ok := pathValue.(map[string]interface{}); ok {
										for method, methodValue := range pathItem {
											if methodOp, ok := methodValue.(map[string]interface{}); ok {
												methodOp["security"] = []map[string][]string{
													{"BearerAuth": {}},
												}
												pathItem[method] = methodOp
											}
										}
										paths[pathKey] = pathItem
									}
								}
								merged["paths"] = paths
							}
						}
						if definitions, ok := driverDoc["definitions"].(map[string]interface{}); ok {
							if mergedDefinitions, ok := merged["definitions"].(map[string]interface{}); ok {
								for k, v := range definitions {
									mergedDefinitions[k] = v
								}
							} else {
								merged["definitions"] = definitions
							}
						}
					}
				}
			} else {
				log.Printf("Driver service config not found")
			}

			c.JSON(http.StatusOK, merged)

		case "/swagger-initializer.js":
			c.File("docs/swagger-initializer.js")

		default:
			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		}
	}
}
