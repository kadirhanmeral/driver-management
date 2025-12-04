package repository

import (
	"context"

	"time"

	"github.com/kadirhanmeral/driver-management/internal/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DriverRepository interface
type DriverRepository interface {
	Create(driver entities.Driver, ctx context.Context) (primitive.ObjectID, error)
	FindByParams(minLat, maxLat, minLon, maxLon *float64, taxiType *string, page, pageSize *int, ctx context.Context) ([]*entities.Driver, error)
	GetByID(id primitive.ObjectID, ctx context.Context) (*entities.Driver, error)
	Update(id primitive.ObjectID, update bson.M, ctx context.Context) error
	Delete(id primitive.ObjectID, ctx context.Context) (int64, error)
}

// MongoDriverRepository struct
type mongoDriverRepository struct {
	collection *mongo.Collection
}

// NewDriverRepository creates a new repository
func NewDriverRepository(client *mongo.Client, dbName, collectionName string) DriverRepository {
	return &mongoDriverRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

// Create adds a new driver
func (r *mongoDriverRepository) Create(driver entities.Driver, ctx context.Context) (primitive.ObjectID, error) {
	driver.CreatedAt = time.Now().UTC()
	driver.UpdatedAt = time.Now().UTC()
	result, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

// FindByParams finds drivers based on location and taxi type with pagination
func (r *mongoDriverRepository) FindByParams(
	minLat, maxLat, minLon, maxLon *float64,
	taxiType *string,
	page, pageSize *int,
	ctx context.Context,
) ([]*entities.Driver, error) {

	filter := bson.M{}

	// Lat filters
	if minLat != nil || maxLat != nil {
		latFilter := bson.M{}
		if minLat != nil {
			latFilter["$gte"] = *minLat
		}
		if maxLat != nil {
			latFilter["$lte"] = *maxLat
		}
		filter["location.lat"] = latFilter
	}

	// Lon filters
	if minLon != nil || maxLon != nil {
		lonFilter := bson.M{}
		if minLon != nil {
			lonFilter["$gte"] = *minLon
		}
		if maxLon != nil {
			lonFilter["$lte"] = *maxLon
		}
		filter["location.lon"] = lonFilter
	}

	// TaxiType filter
	if taxiType != nil && *taxiType != "" {
		filter["taxiType"] = *taxiType
	}

	// Projection
	projection := bson.M{
		"firstName": 1,
		"lastName":  1,
		"plate":     1,
		"location":  1,
		"taxiType":  1,
		"carBrand":  1,
		"carModel":  1,
	}

	findOptions := options.Find()

	// Projection set
	findOptions.SetProjection(projection)

	// Pagination varsa uygula
	if page != nil && pageSize != nil {
		findOptions.SetSkip(int64((*page - 1) * *pageSize))
		findOptions.SetLimit(int64(*pageSize))
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []*entities.Driver
	for cursor.Next(ctx) {
		var d entities.Driver
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		drivers = append(drivers, &d)
	}

	return drivers, nil
}

// GetByID retrieves a single driver by ID
func (r *mongoDriverRepository) GetByID(id primitive.ObjectID, ctx context.Context) (*entities.Driver, error) {
	var driver entities.Driver
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

// Update modifies an existing driver
func (r *mongoDriverRepository) Update(id primitive.ObjectID, update bson.M, ctx context.Context) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Delete removes a driver
func (r *mongoDriverRepository) Delete(id primitive.ObjectID, ctx context.Context) (int64, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
