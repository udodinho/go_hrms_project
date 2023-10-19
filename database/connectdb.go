package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

type MongoInstance struct {
	Client *mongo.Client
	Db	   *mongo.Database
}

var MG MongoInstance

func Connect() error{
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
	
	 dbURI := os.Getenv("MONGO_URI")

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbURI).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
        log.Fatal("Could not connect to mongoDB: ", err)
    }
	
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	
	db := client.Database("GO-HRMS")

	MG = MongoInstance{
		Client: client,
		Db: db,
	}

	
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!", result)

	return nil

}
