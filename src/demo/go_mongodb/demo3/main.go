package main

import (
	"context"
	"fmt"
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
	JobName   string    `bson:"jobName"`   // task name
	Command   string    `bson:"command"`   // shell command
	Err       string    `bson:"err"`       // error info
	Content   string    `bson:"content"`   // output of shell command
	TimePoint TimePoint `bson:"timePoint"` // execute time point
}

func main() {
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		record     *LogRecord
		logArr     []interface{} //  log address
		result     *mongo.InsertManyResult
		insertId   interface{} //  objectId
		docId      objectid.ObjectID
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), "mongodb://36.111.184.221:27017", clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose database cron (init before)
	database = client.Database("cron")

	// 3, choose collection log (init before)
	collection = database.Collection("log")

	// 4, insert record(bson)
	record = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	// 5. batch insert multiple documents
	logArr = []interface{}{record, record, record}

	// start insertion
	if result, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println(err)
		return
	}

	// snowflake id: current time + machine ID + current time incrementing counter
	for _, insertId = range result.InsertedIDs {
		// hold the interface{}, reflect to objectID
		docId = insertId.(objectid.ObjectID)
		fmt.Println("auto-incrementing id:", docId.Hex())
	}
}
