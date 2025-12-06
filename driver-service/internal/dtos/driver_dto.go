package dtos

import (
	"github.com/kadirhanmeral/driver-management/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateDriverDTO struct {
	FirstName string  `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string  `json:"lastName" binding:"required,min=2,max=50"`
	Plate     string  `json:"plate" binding:"required,len=8"`
	TaxiType  string  `json:"taxiType" binding:"required,oneof=sari siyah"`
	CarBrand  string  `json:"carBrand" binding:"required,min=2,max=30"`
	CarModel  string  `json:"carModel" binding:"required,min=1,max=30"`
	Lat       float64 `json:"lat" binding:"required,gte=-90,lte=90"`
	Lon       float64 `json:"lon" binding:"required,gte=-180,lte=180"`
}

type UpdateDriverDTO struct {
	FirstName *string  `json:"firstName,omitempty" binding:"omitempty,min=2,max=50"`
	LastName  *string  `json:"lastName,omitempty" binding:"omitempty,min=2,max=50"`
	Plate     *string  `json:"plate,omitempty" binding:"omitempty,len=8"`
	TaxiType  *string  `json:"taxiType,omitempty" binding:"omitempty,oneof=sari siyah"`
	CarBrand  *string  `json:"carBrand,omitempty" binding:"omitempty,min=2,max=30"`
	CarModel  *string  `json:"carModel,omitempty" binding:"omitempty,min=1,max=30"`
	Lat       *float64 `json:"lat,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Lon       *float64 `json:"lon,omitempty" binding:"omitempty,gte=-180,lte=180"`
}

type DriverResponseNearbyDTO struct {
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Plate      string   `json:"plate"`
	CarBrand   string   `json:"carBrand"`
	DistanceKm *float64 `json:"distanceKm,omitempty"`
}

type DriverResponseDTO struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Plate     string  `json:"plate"`
	TaxiType  string  `json:"taxiType"`
	CarBrand  string  `json:"carBrand"`
	CarModel  string  `json:"carModel"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
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

func (dto *UpdateDriverDTO) ToBsonMap() bson.M {
	update := bson.M{}

	if dto.FirstName != nil {
		update["firstName"] = *dto.FirstName
	}
	if dto.LastName != nil {
		update["lastName"] = *dto.LastName
	}
	if dto.Plate != nil {
		update["plate"] = *dto.Plate
	}
	if dto.TaxiType != nil {
		update["taxiType"] = *dto.TaxiType
	}
	if dto.CarBrand != nil {
		update["carBrand"] = *dto.CarBrand
	}
	if dto.CarModel != nil {
		update["carModel"] = *dto.CarModel
	}
	if dto.Lat != nil {
		update["location.lat"] = *dto.Lat
	}
	if dto.Lon != nil {
		update["location.lon"] = *dto.Lon
	}

	return update
}

func DriverEntityToDriverResponseDTO(driver *entities.Driver) *DriverResponseDTO {
	return &DriverResponseDTO{
		ID:        driver.ID,
		FirstName: driver.FirstName,
		LastName:  driver.LastName,
		Plate:     driver.Plate,
		TaxiType:  driver.TaxiType,
		CarBrand:  driver.CarBrand,
		CarModel:  driver.CarModel,
		Lat:       driver.Location.Lat,
		Lon:       driver.Location.Lon,
	}
}
