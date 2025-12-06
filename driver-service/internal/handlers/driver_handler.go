package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kadirhanmeral/driver-management/internal/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"

	services "github.com/kadirhanmeral/driver-management/internal/services"
)

type DriverHandler struct {
	service *services.DriverService
}

func NewDriverHandler(service *services.DriverService) *DriverHandler {
	return &DriverHandler{service: service}
}

// CreateDriver godoc
// @Summary      Create a new driver
// @Description  Create a new driver with the input payload
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        driver  body      dtos.CreateDriverDTO  true  "Create Driver"
// @Success      201     {object}  map[string]string
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /drivers [post]
func (h *DriverHandler) CreateDriver(ctx *gin.Context) {

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

// Get driver by ID
// GetDriver godoc
// @Summary      Get a driver
// @Description  Get a driver by ID
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Driver ID"
// @Success      200  {object}  dtos.DriverResponseDTO
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /drivers/{id} [get]
func (h *DriverHandler) GetDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	driver, err := h.service.GetDriver(objID, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

// Update driver by ID and body
// UpdateDriver godoc
// @Summary      Update a driver
// @Description  Update a driver by ID
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        id      path      string                true  "Driver ID"
// @Param        driver  body      dtos.UpdateDriverDTO  true  "Update Driver"
// @Success      204     "No Content"
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /drivers/{id} [patch]
func (h *DriverHandler) UpdateDriver(ctx *gin.Context) {

	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	var updateDto dtos.UpdateDriverDTO
	if err := ctx.ShouldBindJSON(&updateDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateMap := updateDto.ToBsonMap()
	if len(updateMap) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}
	updateMap["updatedAt"] = time.Now().UTC()

	if err := h.service.UpdateDriver(objID, updateMap, ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// Delete driver with ID
// DeleteDriver godoc
// @Summary      Delete a driver
// @Description  Delete a driver by ID
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Driver ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /drivers/{id} [delete]
func (h *DriverHandler) DeleteDriver(ctx *gin.Context) {
	driverID := ctx.Param("id")
	objID, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	deletedCount, err := h.service.DeleteDriver(objID, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if deletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Driver not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Driver deleted"})
}

// ListDrivers godoc
// @Summary      List drivers
// @Description  Get all drivers with pagination
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        pageSize  query     int     false  "Page size"
// @Success      200       {array}   dtos.DriverResponseDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /drivers [get]
func (h *DriverHandler) ListDrivers(ctx *gin.Context) {
	var pagePtr, pageSizePtr *int

	if pageStr := ctx.Query("page"); pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil {
			if v < 1 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "page must be greater than 0"})
				return
			}
			pagePtr = &v
		}
	}

	if sizeStr := ctx.Query("pageSize"); sizeStr != "" {
		if v, err := strconv.Atoi(sizeStr); err == nil {
			if v < 1 || v > 100 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "pageSize must be between 1 and 100"})
				return
			}
			pageSizePtr = &v
		}
	}

	drivers, err := h.service.ListDrivers(pagePtr, pageSizePtr, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, drivers)
}

// GetNearbyDrivers godoc
// @Summary      List nearby drivers
// @Description  Get drivers within 6 km of a location
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Param        lat       query     number  true   "Latitude"
// @Param        lon       query     number  true   "Longitude"
// @Param        taxiType  query     string  false  "Taxi Type"
// @Success      200       {array}   dtos.DriverResponseNearbyDTO
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /drivers/nearby [get]
func (h *DriverHandler) GetNearbyDrivers(ctx *gin.Context) {
	latStr := ctx.Query("lat")
	lonStr := ctx.Query("lon")

	if latStr == "" || lonStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lat and lon are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lat must be number"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "lon must be number"})
		return
	}

	taxiType := ctx.Query("taxiType")

	drivers, err := h.service.GetNearbyDrivers(lat, lon, taxiType, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, drivers)
}
