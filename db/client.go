package db

import (
	"context"
	"time"

	"github.com/cicovic-andrija/spot/resources"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client represents a database client object
type Client struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewClient(connstring, database, collection string) (*Client, error) {
	ctx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connstring))
	if err != nil {
		return nil, err
	}

	ctx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Client{client: client, database: database, collection: collection}, nil
}

func (c *Client) FindAllGarages(ctx context.Context) (map[string]*resources.Garage, error) {
	collection := c.client.Database(c.database).Collection(c.collection)

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	garages := make(map[string]*resources.Garage)
	for cursor.Next(ctx) {
		g := &resources.Garage{}
		err = cursor.Decode(g)
		if err != nil {
			return nil, err
		}
		garages[g.ID] = g
	}

	return garages, nil
}

func (c *Client) InsertGarage(ctx context.Context, garage *resources.Garage) error {
	collection := c.client.Database(c.database).Collection(c.collection)

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"id":          garage.ID,
			"name":        garage.Name,
			"city":        garage.City,
			"address":     garage.Address,
			"geolocation": garage.Geolocation,
			"sections":    garage.Sections,
		},
	)
	return err
}

func (c *Client) UpdateGarage(ctx context.Context, id string, newname string, newcity string, newaddress string) error {
	collection := c.client.Database(c.database).Collection(c.collection)

	update := bson.M{}
	if newname != "" {
		update["name"] = newname
	}
	if newcity != "" {
		update["city"] = newcity
	}
	if newaddress != "" {
		update["address"] = newaddress
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"id": id,
		},
		bson.M{
			"$set": update,
		},
	)
	return err
}

func (c *Client) DeleteGarage(ctx context.Context, id string) error {
	collection := c.client.Database(c.database).Collection(c.collection)
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (c *Client) InsertSection(ctx context.Context, garageID string, section *resources.Section) error {
	collection := c.client.Database(c.database).Collection(c.collection)

	update := bson.M{
		"sections": bson.M{
			"name":        section.Name,
			"level":       section.Level,
			"description": section.Description,
			"total_spots": section.TotalSpots,
		},
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"id": garageID,
		},
		bson.M{
			"$push": update,
		},
	)
	return err
}

func (c *Client) UpdateSection(
	ctx context.Context,
	garageID string,
	sectionName string,
	newname string,
	newlevel string,
	newdescription string,
	newtotalspots int,
) error {
	collection := c.client.Database(c.database).Collection(c.collection)

	update := bson.M{}
	if newname != "" {
		update["sections.$.name"] = newname
	}
	if newlevel != "" {
		update["sections.$.level"] = newlevel
	}
	if newdescription != "" {
		update["sections.$.description"] = newdescription
	}
	if newtotalspots > 0 {
		update["sections.$.total_spots"] = newtotalspots
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"id":            garageID,
			"sections.name": sectionName,
		},
		bson.M{
			"$set": update,
		},
	)
	return err
}

func (c *Client) DeleteSection(ctx context.Context, garageID string, sectionName string) error {
	collection := c.client.Database(c.database).Collection(c.collection)

	update := bson.M{
		"sections": bson.M{
			"name": sectionName,
		},
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{
			"id": garageID,
		},
		bson.M{
			"$pull": update,
		},
	)
	return err
}
