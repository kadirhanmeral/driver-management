package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/internal/handlers"

	_ "github.com/kadirhanmeral/driver-management/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterDriverEndpoints(router *gin.Engine, driverHandlers *handlers.DriverHandler) {

	router.POST("/drivers", driverHandlers.CreateDriver)
	router.GET("/drivers", driverHandlers.ListDrivers)
	router.GET("/drivers/nearby", driverHandlers.GetNearbyDrivers)
	router.GET("/drivers/:id", driverHandlers.GetDriver)
	router.PATCH("/drivers/:id", driverHandlers.UpdateDriver)
	router.DELETE("/drivers/:id", driverHandlers.DeleteDriver)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
