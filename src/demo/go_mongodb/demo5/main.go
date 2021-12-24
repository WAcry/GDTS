package main

import (
	"context"
	"demo/config"
	"fmt"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
)

// startTime < some time
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp} }
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
		delResult  *mongo.DeleteResult
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), config.MONGODB_URL, clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose database cron (init before)
	database = client.Database("cron")

	// 3, choose collection log (init before)
	collection = database.Collection("log")

	// 4, delete all logs which startTime < now	//($lt = less than)
	//  delete({"timePoint.startTime": {"$lt": current time}})
	delCond = &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	// execute deletion
	if delResult, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("delete ", delResult.DeletedCount, " rows")
}
