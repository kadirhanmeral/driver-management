package entities

import "time"

type Driver struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	FirstName string    `bson:"firstName" json:"firstName"`
	LastName  string    `bson:"lastName" json:"lastName"`
	Plate     string    `bson:"plate" json:"plate"`
	TaxiType  string    `bson:"taxiType" json:"taxiType"`
	CarBrand  string    `bson:"carBrand" json:"carBrand"`
	CarModel  string    `bson:"carModel" json:"carModel"`
	Location  GeoPoint  `bson:"location" json:"location"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type GeoPoint struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lon float64 `bson:"lon" json:"lon"`
}
