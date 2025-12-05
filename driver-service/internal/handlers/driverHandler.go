package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/internal/dtos"
	services "github.com/kadirhanmeral/driver-management/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Driver struct {
	service *services.DriverService
}

func NewDriverHandler(service *services.DriverService) *Driver {
	return &Driver{service: service}
}

// --------------------
// Tüm driver’ları listele (pagination destekli)
// GET /drivers?page=1&pageSize=20
// --------------------
func (h *Driver) ListDrivers(ctx *gin.Context) {
	var pagePtr *int
	var pageSizePtr *int

	// page var mı?
	if pageStr, exists := ctx.GetQuery("page"); exists {
		pageVal, err := strconv.Atoi(pageStr)
		if err == nil {
			pagePtr = &pageVal
		}
	}

	// pageSize var mı?
	if pageSizeStr, exists := ctx.GetQuery("pageSize"); exists {
		pageSizeVal, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSizePtr = &pageSizeVal
		}
	}

	drivers, err := h.service.ListDrivers(pagePtr, pageSizePtr, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, drivers)
}

// --------------------
// Yakındaki taksileri bul
// GET /drivers/nearby?lat=41.001&lon=28.99&taxiType=sari&page=1&pageSize=20
// --------------------
func (h *Driver) GetNearbyDrivers(ctx *gin.Context) {
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")
	taxiType := ctx.Query("taxiType")

	if latStr == "" || lonStr == "" || taxiType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lat, lon ve taxiType parametreleri gerekli"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lat parametresi geçerli bir sayı olmalı"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lon parametresi geçerli bir sayı olmalı"})
		return
	}

	drivers, err := h.service.GetNearbyDrivers(lat, lon, taxiType, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, drivers)
}

// --------------------
// Yeni driver ekle
// POST /drivers
// --------------------
func (h *Driver) CreateDriver(ctx *gin.Context) {
	var createDto dtos.CreateDriverDTO
	if err := ctx.ShouldBindJSON(&createDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	driver := createDto.ToEntity()
	id, err := h.service.CreateDriver(driver, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

// --------------------
// Tek driver getir
// GET /drivers/:id
// --------------------
func (h *Driver) GetDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz driver ID"})
		return
	}

	driver, err := h.service.GetDriver(objID, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver bulunamadı"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

// --------------------
// Driver güncelle
// PUT /drivers/:id
// --------------------
func (h *Driver) UpdateDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz driver ID"})
		return
	}

	var updateDto dtos.UpdateDriverDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mevcut driver'ı al
	driver, err := h.service.GetDriver(objID, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver bulunamadı"})
		return
	}

	updated := updateDto.ToEntity(driver)

	err = h.service.UpdateDriver(objID, map[string]interface{}{
		"firstName": updated.FirstName,
		"lastName":  updated.LastName,
		"plate":     updated.Plate,
		"taxiType":  updated.TaxiType,
		"carBrand":  updated.CarBrand,
		"carModel":  updated.CarModel,
		"location":  updated.Location,
	}, ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Driver güncellendi"})
}

// --------------------
// Driver sil
// DELETE /drivers/:id
// --------------------
func (h *Driver) DeleteDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz driver ID"})
		return
	}

	deletedCount, err := h.service.DeleteDriver(objID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Driver bulunamadı"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Driver silindi"})
}
