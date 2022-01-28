package main

import (
	"context"
	"fmt"
	"lib/config"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
)

func main() {
	var (
		client     *mongo.Client
		err        error
		db         *mongo.Database
		collection *mongo.Collection
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), config.MONGODB_URL, clientopt.ConnectTimeout(30*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose database my_db (init before)
	db = client.Database("my_db")

	// 3, choose collection my_collection (init before)
	collection = db.Collection("my_collection")

	collection = collection
}
