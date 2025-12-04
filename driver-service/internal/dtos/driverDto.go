package dtos

import (
	"time"

	"github.com/kadirhanmeral/driver-management/internal/entities"
)

type CreateDriverDTO struct {
	FirstName string  `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string  `json:"lastName" binding:"required,min=2,max=50"`
	Plate     string  `json:"plate" binding:"required,len=8"`
	TaxiType  string  `json:"taxiType" binding:"required,oneof=sari beyaz"`
	CarBrand  string  `json:"carBrand" binding:"required,min=2,max=30"`
	CarModel  string  `json:"carModel" binding:"required,min=1,max=30"`
	Lat       float64 `json:"lat" binding:"required,gte=-90,lte=90"`
	Lon       float64 `json:"lon" binding:"required,gte=-180,lte=180"`
}
type UpdateDriverDTO struct {
	FirstName *string  `json:"firstName,omitempty" binding:"omitempty,min=2,max=50"`
	LastName  *string  `json:"lastName,omitempty" binding:"omitempty,min=2,max=50"`
	Plate     *string  `json:"plate,omitempty" binding:"omitempty,len=8"`
	TaxiType  *string  `json:"taxiType,omitempty" binding:"omitempty,oneof=sari beyaz"`
	CarBrand  *string  `json:"carBrand,omitempty" binding:"omitempty,min=2,max=30"`
	CarModel  *string  `json:"carModel,omitempty" binding:"omitempty,min=1,max=30"`
	Lat       *float64 `json:"lat,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Lon       *float64 `json:"lon,omitempty" binding:"omitempty,gte=-180,lte=180"`
}

type DriverResponseDTO struct {
	ID         string            `json:"id"`
	FirstName  string            `json:"firstName"`
	LastName   string            `json:"lastName"`
	Plate      string            `json:"plate"`
	TaxiType   string            `json:"taxiType"`
	CarBrand   string            `json:"carBrand"`
	CarModel   string            `json:"carModel"`
	Location   entities.GeoPoint `json:"location"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	DistanceKm *float64          `json:"distanceKm,omitempty"` // nearby endpoint için
}

func (dto *CreateDriverDTO) ToEntity() *entities.Driver {
	return &entities.Driver{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Plate:     dto.Plate,
		TaxiType:  dto.TaxiType,
		CarBrand:  dto.CarBrand,
		CarModel:  dto.CarModel,
		Location: entities.GeoPoint{
			Lat: dto.Lat,
			Lon: dto.Lon,
		},
	}
}

func (dto *UpdateDriverDTO) ToEntity(driver *entities.Driver) *entities.Driver {
	if dto.FirstName != nil {
		driver.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		driver.LastName = *dto.LastName
	}
	if dto.Plate != nil {
		driver.Plate = *dto.Plate
	}
	if dto.TaxiType != nil {
		driver.TaxiType = *dto.TaxiType
	}
	if dto.CarBrand != nil {
		driver.CarBrand = *dto.CarBrand
	}
	if dto.CarModel != nil {
		driver.CarModel = *dto.CarModel
	}
	if dto.Lat != nil {
		driver.Location.Lat = *dto.Lat
	}
	if dto.Lon != nil {
		driver.Location.Lon = *dto.Lon
	}
	return driver
}

func (dto *DriverResponseDTO) FromEntity(driver *entities.Driver) *DriverResponseDTO {
	return &DriverResponseDTO{
		ID:        driver.ID,
		FirstName: driver.FirstName,
		LastName:  driver.LastName,
		Plate:     driver.Plate,
		TaxiType:  driver.TaxiType,
		CarBrand:  driver.CarBrand,
		CarModel:  driver.CarModel,
		Location:  driver.Location,
		CreatedAt: driver.CreatedAt,
		UpdatedAt: driver.UpdatedAt,
	}
}
