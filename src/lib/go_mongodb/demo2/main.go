package main

import (
	"context"
	"fmt"
	"lib/config"
	"time"
	"utils/mongodb/mongo-go-driver/bson/objectid"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
)

// TimePoint time point of job execution
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// LogRecord record of log
type LogRecord struct {
	JobName   string    `bson:"jobName"`   // job name
	Command   string    `bson:"command"`   // shell command
	Err       string    `bson:"err"`       // error info
	Content   string    `bson:"content"`   // output of shell command
	TimePoint TimePoint `bson:"timePoint"` // execute time point
}

func main() {
	var (
		client     *mongo.Client
		err        error
		db         *mongo.Database
		collection *mongo.Collection
		rec        *LogRecord
		result     *mongo.InsertOneResult
		docId      objectid.ObjectID
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), config.MONGODB_URL, clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose db gdts (init before)
	db = client.Database("gdts")

	// 3, choose collection log (init before)
	collection = db.Collection("log")

	// 4, insert rec(bson)
	rec = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	if result, err = collection.InsertOne(context.TODO(), rec); err != nil {
		fmt.Println(err)
		return
	}

	// _id: generate a globally unique ObjectID, 12-byte binary.
	docId = result.InsertedID.(objectid.ObjectID)
	fmt.Println("auto-incrementing id:", docId.Hex())
}
