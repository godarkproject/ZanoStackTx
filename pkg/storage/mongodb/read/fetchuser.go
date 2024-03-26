package mongodb

import (
	_ "compress/zlib"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type ZanoDeposits struct {
	TxHash string `bson:"tx_hash"`
	Amount int64  `bson:"amount"`
}

type GameTransactions struct {
	GameUuid string `bson:"game_uuid"`
	Amount   int64  `bson:"amount"`
}

type Withdrawals struct {
	TxHash string `bson:"tx_hash"`
	Amount int64  `bson:"amount"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	CreationTime     int64              `bson:"creation_time"`
	Ip               string             `bson:"ip"`
	Uuid             string             `bson:"uuid"`
	Username         string             `bson:"username"`
	Password         string             `bson:"password"`
	ZanoAddress      string             `bson:"zano_address"`
	ZanoPaymentId    string             `bson:"zano_payment_id"`
	ExternalAddress  string             `bson:"external_address"`
	ZanoDeposits     []ZanoDeposits     `bson:"zano_deposits"`
	GameTransactions []GameTransactions `bson:"game_transactions"`
	Withdrawals      []Withdrawals      `bson:"withdrawals"`
	Balance          int64              `bson:"balance"`
}

func FetchUser(mongoUri string, paymentId string) (User, error) {
	fmt.Printf("fetching user")
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
	// Creates a query filter to match documents in which the "name" is
	// "Bagels N Buns"
	filter := bson.D{{"zano_payment_id", paymentId}}
	// Retrieves the first matching document
	var result User
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	// Prints a message if no documents are matched or if any
	// other errors occur during the operation
	if err != nil {
		if errors.Is(mongo.ErrNoDocuments, err) {
			return result, err
		}
		panic(err)
	}

	return result, nil
}
