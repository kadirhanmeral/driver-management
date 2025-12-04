package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/internal/handlers"
)

// RegisterDriverEndpoints driver ile ilgili tüm endpoint'leri kaydeder
func RegisterDriverEndpoints(router *gin.Engine, driverHandlers *handlers.Driver) {
	// CRUD endpoints
	router.POST("/drivers", driverHandlers.CreateDriver)       // Yeni driver ekle
	router.GET("/drivers", driverHandlers.ListDrivers)         // Tüm driver'ları listele
	router.GET("/drivers/:id", driverHandlers.GetDriver)       // Tek driver bilgisi
	router.PUT("/drivers/:id", driverHandlers.UpdateDriver)    // Driver güncelle
	router.DELETE("/drivers/:id", driverHandlers.DeleteDriver) // Driver sil

	// Yakındaki taksileri listele
	router.GET("/drivers/nearby", driverHandlers.GetNearbyDrivers) // Query param: lat, lon, taxiType
}
