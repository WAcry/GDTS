package main

import (
	"context"
	"demo/config"
	"fmt"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), config.MONGODB_URL, clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose database my_db (init before)
	database = client.Database("my_db")

	// 3, choose collection my_collection (init before)
	collection = database.Collection("my_collection")

	collection = collection
}
