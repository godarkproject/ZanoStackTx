package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func UpdateBalance(mongoUri string, balance int64, mongoId primitive.ObjectID) (bool, error) {
	opts := options.Client().ApplyURI(fmt.Sprintf("%s?compressors=snappy,zlib,zstd", mongoUri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// If you wish to know if a MongoDB server has been found and connected to, use the Ping method:
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("MongoDB server has not been found")
	}

	coll := client.Database("zanostack").Collection("users")
	id, _ := primitive.ObjectIDFromHex(mongoId.Hex())
	filter := bson.D{{"_id", id}}
	// Creates instructions to add the "avg_rating" field to documents
	update := bson.D{{"$set", bson.D{{"balance", balance}}}}
	// Updates the first document that has the specified "_id" value
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}

	return true, nil

}
