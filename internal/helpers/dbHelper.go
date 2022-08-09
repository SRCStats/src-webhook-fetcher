package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	s "github.com/suprnova/src-webhook-fetcher/internal/structs"
)

var (
	database   string
	collection string
)

const (
	connectionStringEnv = "SRC_WEBHOOK_MONGODB_CONNECTION_STRING"
	databaseNameEnv     = "SRC_WEBHOOK_DATABASE"
	collectionNameEnv   = "SRC_WEBHOOK_COLLECTION"
)

func Connect() *mongo.Client {
	connectionString := os.Getenv(connectionStringEnv)
	if connectionString == "" {
		log.Fatal("The database connection string variable is missing!")
	}
	database = os.Getenv(databaseNameEnv)
	if database == "" {
		log.Fatal("The database name variable is missing!")
	}
	collection = os.Getenv(collectionNameEnv)
	if collection == "" {
		log.Fatal("The collection name variable is missing!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	options := options.Client().ApplyURI(connectionString).SetDirect(true)
	c, _ := mongo.NewClient(options)

	err := c.Connect(ctx)
	if err != nil {
		log.Fatalf("Connection couldn't be initialized!\n%v", err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Connection didn't respond back!\n%v", err)
	}
	return c
}

func Create(runs s.Response) {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	r, err := runCollection.InsertOne(ctx, runs)
	if err != nil {
		log.Fatalf("The runs couldn't be added!\n%v", err)
	}
	fmt.Println("Added runs", r.InsertedID)
}

func List(scope string) *s.Response {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)
	var runs s.Response
	runCollection := c.Database(database).Collection(collection)
	err := runCollection.FindOne(ctx, bson.D{{Key: "scope", Value: scope}}).Decode(&runs)

	if err != nil {
		log.Printf("The runs couldn't be listed!\n%v", err)
	}
	return &runs
}

func Delete(scope string) {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	res, err := runCollection.DeleteMany(ctx, bson.D{{Key: "scope", Value: scope}})
	if err != nil {
		log.Fatalf("The runs couldn't be deleted!\n%v", err)
	}
	fmt.Printf("Deleted %v entry for scope %v\n", res.DeletedCount, scope)
}
