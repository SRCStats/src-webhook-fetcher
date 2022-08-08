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

func Create(runs *[]s.Run) {
	runsItf := make([]interface{}, len(*runs))
	for i, run := range *runs {
		runsItf[i] = run
	}
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	r, err := runCollection.InsertMany(ctx, runsItf)
	if err != nil {
		log.Fatalf("The runs couldn't be added!\n%v", err)
	}
	fmt.Println("Added runs", r.InsertedIDs)
}

func List(scope string) *[]s.Run {
	c := Connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	runCollection := c.Database(database).Collection(collection)
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	rs, err := runCollection.Find(ctx, bson.D{{Key: "scope", Value: scope}}, opts)
	var runs []s.Run

	if err != nil {
		log.Fatalf("The runs couldn't be listed!\n%v", err)
	}
	err = rs.All(ctx, &runs)
	if err != nil {
		log.Fatalf("The runs couldn't be listed!\n%v", err)
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
	fmt.Printf("Deleted %v runs\n", res.DeletedCount)
}
