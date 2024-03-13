package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func UpdateBalance(balance int64, mongoId primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

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
