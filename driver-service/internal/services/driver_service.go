package service

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/kadirhanmeral/driver-management/internal/dtos"
	"github.com/kadirhanmeral/driver-management/internal/entities"
	"github.com/kadirhanmeral/driver-management/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DriverService struct {
	repo repository.DriverRepository
}

func NewDriverService(repo repository.DriverRepository) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) CreateDriver(driver *entities.Driver, ctx context.Context) (primitive.ObjectID, error) {
	return s.repo.Create(*driver, ctx)
}

func (s *DriverService) GetDriver(id primitive.ObjectID, ctx context.Context) (*dtos.DriverResponseDTO, error) {

	driverEntity, err := s.repo.GetByID(id, ctx)

	if err != nil {
		return nil, err
	}

	return dtos.DriverEntityToDriverResponseDTO(driverEntity), nil
}

func (s *DriverService) UpdateDriver(id primitive.ObjectID, update map[string]interface{}, ctx context.Context) error {
	update["updatedAt"] = time.Now().UTC()
	return s.repo.Update(id, update, ctx)
}

func (s *DriverService) DeleteDriver(id primitive.ObjectID, ctx context.Context) (int64, error) {
	return s.repo.Delete(id, ctx)
}

func (s *DriverService) ListDrivers(page, pageSize *int, ctx context.Context) ([]*dtos.DriverResponseDTO, error) {
	driverEntityList, err := s.repo.FindByParams(page, pageSize, ctx)
	if err != nil {
		return nil, err
	}

	driverResponseDTOList := make([]*dtos.DriverResponseDTO, 0)

	for _, driverEntity := range driverEntityList {
		dto := dtos.DriverEntityToDriverResponseDTO(driverEntity)
		driverResponseDTOList = append(driverResponseDTOList, dto)
	}

	return driverResponseDTOList, nil
}

func (s *DriverService) GetNearbyDrivers(lat, lon float64, taxiType string, ctx context.Context) ([]*dtos.DriverResponseNearbyDTO, error) {
	const earthRadius = 6371.0
	const maxDistanceKm = 6.0

	latDiff := (maxDistanceKm / earthRadius) * (180 / math.Pi)
	lonDiff := (maxDistanceKm / earthRadius) * (180 / math.Pi) / math.Cos(lat*math.Pi/180)

	minLat := lat - latDiff
	maxLat := lat + latDiff
	minLon := lon - lonDiff
	maxLon := lon + lonDiff

	minLatPtr := &minLat
	maxLatPtr := &maxLat
	minLonPtr := &minLon
	maxLonPtr := &maxLon
	taxiTypePtr := &taxiType

	drivers, err := s.repo.FindByParamsNearby(minLatPtr, maxLatPtr, minLonPtr, maxLonPtr, taxiTypePtr, ctx)
	if err != nil {
		return nil, err
	}

	var nearby []*dtos.DriverResponseNearbyDTO
	for _, d := range drivers {

		distance := Haversine(lat, lon, d.Location.Lat, d.Location.Lon)

		if distance <= maxDistanceKm {
			dto := &dtos.DriverResponseNearbyDTO{
				FirstName:  d.FirstName,
				LastName:   d.LastName,
				Plate:      d.Plate,
				DistanceKm: &distance,
			}
			nearby = append(nearby, dto)
		}
	}

	sort.Slice(nearby, func(i, j int) bool {
		return *nearby[i].DistanceKm < *nearby[j].DistanceKm
	})

	return nearby, nil
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius in kilometers

	dLat := (lat2 - lat1) * math.Pi / 180.0 // Δφ
	dLon := (lon2 - lon1) * math.Pi / 180.0 // Δλ

	lat1Rad := lat1 * math.Pi / 180.0 // φ1
	lat2Rad := lat2 * math.Pi / 180.0 // φ2

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
