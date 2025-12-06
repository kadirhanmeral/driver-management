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

type DriverRepository interface {
	Create(driver entities.Driver, ctx context.Context) (primitive.ObjectID, error)
	FindByParamsNearby(minLat, maxLat, minLon, maxLon *float64, taxiType *string, ctx context.Context) ([]*entities.Driver, error)
	FindByParams(page, pageSize *int, ctx context.Context) ([]*entities.Driver, error)
	GetByID(id primitive.ObjectID, ctx context.Context) (*entities.Driver, error)
	Update(id primitive.ObjectID, update bson.M, ctx context.Context) error
	Delete(id primitive.ObjectID, ctx context.Context) (int64, error)
}

type mongoDriverRepository struct {
	collection *mongo.Collection
}

func NewDriverRepository(client *mongo.Client, dbName, collectionName string) DriverRepository {
	return &mongoDriverRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func (r *mongoDriverRepository) Create(driver entities.Driver, ctx context.Context) (primitive.ObjectID, error) {
	driver.CreatedAt = time.Now().UTC()
	driver.UpdatedAt = time.Now().UTC()
	result, err := r.collection.InsertOne(ctx, driver)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *mongoDriverRepository) FindByParamsNearby(
	minLat, maxLat, minLon, maxLon *float64,
	taxiType *string,
	ctx context.Context,
) ([]*entities.Driver, error) {

	filter := bson.M{}

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

	if taxiType != nil && *taxiType != "" {
		filter["taxiType"] = *taxiType
	}

	projection := bson.M{
		"firstName": 1,
		"lastName":  1,
		"plate":     1,
		"location":  1,
	}

	findOptions := options.Find()

	findOptions.SetProjection(projection)

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

func (r *mongoDriverRepository) FindByParams(
	page, pageSize *int,
	ctx context.Context,
) ([]*entities.Driver, error) {

	filter := bson.M{}

	findOptions := options.Find()

	if page != nil && pageSize != nil {
		skip := int64((*page - 1) * *pageSize)
		limit := int64(*pageSize)
		findOptions.SetSkip(skip)
		findOptions.SetLimit(limit)
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

func (r *mongoDriverRepository) GetByID(id primitive.ObjectID, ctx context.Context) (*entities.Driver, error) {
	var driver entities.Driver
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *mongoDriverRepository) Update(id primitive.ObjectID, update bson.M, ctx context.Context) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

func (r *mongoDriverRepository) Delete(id primitive.ObjectID, ctx context.Context) (int64, error) {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
