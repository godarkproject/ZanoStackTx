package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func FetchUser(paymentId string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://kekzploit:machinecodes@localhost:27017/"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

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
