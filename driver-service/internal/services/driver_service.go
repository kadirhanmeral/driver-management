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

// CreateDriver adds a new driver
func (s *DriverService) CreateDriver(driver *entities.Driver, ctx context.Context) (primitive.ObjectID, error) {
	return s.repo.Create(*driver, ctx)
}

// GetDriver retrieves a driver by ID
func (s *DriverService) GetDriver(id primitive.ObjectID, ctx context.Context) (*entities.Driver, error) {
	return s.repo.GetByID(id, ctx)
}

// UpdateDriver updates driver details
func (s *DriverService) UpdateDriver(id primitive.ObjectID, update map[string]interface{}, ctx context.Context) error {
	update["updatedAt"] = time.Now().UTC()
	return s.repo.Update(id, update, ctx)
}

// DeleteDriver removes a driver by ID
func (s *DriverService) DeleteDriver(id primitive.ObjectID, ctx context.Context) (int64, error) {
	return s.repo.Delete(id, ctx)
}

// ListDrivers returns a paginated list of drivers
func (s *DriverService) ListDrivers(page, pageSize *int, ctx context.Context) ([]*entities.Driver, error) {
	return s.repo.FindByParams(page, pageSize, ctx)
}

func (s *DriverService) GetNearbyDrivers(lat, lon float64, taxiType string, ctx context.Context) ([]*dtos.DriverResponseDTO, error) {
	const earthRadius = 6371.0
	const maxDistanceKm = 6.0

	// Calculate Bounding Box
	// We create a square area (bounding box) around the user to filter the initial
	latDiff := (maxDistanceKm / earthRadius) * (180 / math.Pi)
	lonDiff := (maxDistanceKm / earthRadius) * (180 / math.Pi) / math.Cos(lat*math.Pi/180)

	// Define the boundaries
	minLat := lat - latDiff
	maxLat := lat + latDiff
	minLon := lon - lonDiff
	maxLon := lon + lonDiff

	// Prepare pointers for the dynamic query
	minLatPtr := &minLat
	maxLatPtr := &maxLat
	minLonPtr := &minLon
	maxLonPtr := &maxLon
	taxiTypePtr := &taxiType

	// Fetch drivers within the square bounding box.
	// This includes "corner" drivers who are technically further than 6km.
	drivers, err := s.repo.FindByParamsNearby(minLatPtr, maxLatPtr, minLonPtr, maxLonPtr, taxiTypePtr, ctx)
	if err != nil {
		return nil, err
	}

	// The database returned a square we need a circle
	var nearby []*dtos.DriverResponseDTO
	for _, d := range drivers {
		// Calculate exact distance: d = R * c
		distance := Haversine(lat, lon, d.Location.Lat, d.Location.Lon)

		// Filter: Keep only if inside the 6km radius
		if distance <= maxDistanceKm {
			dto := &dtos.DriverResponseDTO{
				FirstName:  d.FirstName,
				LastName:   d.LastName,
				Plate:      d.Plate,
				DistanceKm: &distance,
			}
			nearby = append(nearby, dto)
		}
	}

	// Sort by distance (nearest first)
	sort.Slice(nearby, func(i, j int) bool {
		return *nearby[i].DistanceKm < *nearby[j].DistanceKm
	})

	return nearby, nil
}

// Haversine calculates the great-circle distance between two points on a sphere.
// a = sin²(Δφ/2) + cos φ1 ⋅ cos φ2 ⋅ sin²(Δλ/2)
// c = 2 ⋅ atan2( √a, √(1−a) )
// d = R ⋅ c
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius in kilometers

	// Convert degrees to radians
	dLat := (lat2 - lat1) * math.Pi / 180.0 // Δφ
	dLon := (lon2 - lon1) * math.Pi / 180.0 // Δλ

	lat1Rad := lat1 * math.Pi / 180.0 // φ1
	lat2Rad := lat2 * math.Pi / 180.0 // φ2

	// Apply Haversine formula part 'a'
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)

	// Calculate angular distance 'c' in radians
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Final distance 'd'
	return R * c
}
