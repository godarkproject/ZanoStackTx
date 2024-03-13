package mongodb

import (
	"context"
	mongodb "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/read"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func AddTx(txHash string, amount int64, userId primitive.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("zanostack").Collection("users")
	id, err := primitive.ObjectIDFromHex(userId.Hex())
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"_id", id}}

	doc := mongodb.ZanoDeposits{
		TxHash: txHash,
		Amount: amount,
	}

	//update := bson.D{{"$set", bson.D{{"zano_deposits", ""}}}}
	update := bson.M{"$push": bson.M{"zano_deposits": doc}}

	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
}
