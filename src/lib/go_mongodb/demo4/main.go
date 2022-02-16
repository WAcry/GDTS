package main

import (
	"context"
	"fmt"
	"lib/config"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
	"utils/mongodb/mongo-go-driver/mongo/findopt"
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

// FindByJobName filter condition of jobName
type FindByJobName struct {
	JobName string `bson:"jobName"` // JobName赋值为job10
}

func main() {
	// mongodb reads back .bson, and needs to deserialize to LogRecord object
	var (
		client     *mongo.Client
		err        error
		database   *mongo.Database
		collection *mongo.Collection
		cond       *FindByJobName
		cursor     mongo.Cursor
		record     *LogRecord
	)
	// 1, connect to mongodb
	if client, err = mongo.Connect(context.TODO(), config.MONGODB_URL, clientopt.ConnectTimeout(5*time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	// 2, choose database gdts (init before)
	database = client.Database("gdts")

	// 3, choose collection log (init before)
	collection = database.Collection("log")

	// 4. filter by jobName, find 5 documents that jobName=job10
	cond = &FindByJobName{JobName: "job10"} // {"jobName": "job10"}

	// 5. query (filter + pagination parameters)
	if cursor, err = collection.Find(context.TODO(), cond, findopt.Skip(0), findopt.Limit(2)); err != nil { // pagination parameters (start=0, limit=2)
		fmt.Println(err)
		return
	}

	// defer release cursor
	defer cursor.Close(context.TODO())

	// 6, iterate all result set
	for cursor.Next(context.TODO()) {
		// define a record
		record = &LogRecord{}

		// deserialize bson to record object
		if err = cursor.Decode(record); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(*record)
	}
}
