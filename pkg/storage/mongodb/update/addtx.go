package mongodb

import (
	_ "compress/zlib"
	"context"
	"fmt"
	mongodb "github.com/godarkproject/ZanoStackTx/pkg/storage/mongodb/read"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func AddTx(mongoUri string, txHash string, amount int64, userId primitive.ObjectID) {
	fmt.Printf("in add tx")
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
