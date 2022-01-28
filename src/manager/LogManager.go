package manager

import (
	"common"
	"context"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
	"utils/mongodb/mongo-go-driver/mongo/findopt"
)

// MstLogManager mongodb log manager
type MstLogManager struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	LogManager *MstLogManager
)

func InitLogManager() (err error) {
	var (
		client *mongo.Client
	)

	// connect to mongodb
	if client, err = mongo.Connect(
		context.TODO(),
		Config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(Config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		return
	}

	LogManager = &MstLogManager{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}
	return
}

// ListLog list logs
func (logMgr *MstLogManager) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		filter  *common.JobLogFilter
		logSort *common.SortLogByStartTime
		cursor  mongo.Cursor
		jobLog  *common.JobLog
	)

	// len(logArr)
	logArr = make([]*common.JobLog, 0)

	filter = &common.JobLogFilter{JobName: name}

	// sort by start time
	logSort = &common.SortLogByStartTime{SortOrder: -1}

	// query
	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, findopt.Sort(logSort), findopt.Skip(int64(skip)), findopt.Limit(int64(limit))); err != nil {
		return
	}
	// release cursor after done
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		// decode bson
		if err = cursor.Decode(jobLog); err != nil {
			continue // illegal bson, skip
		}

		logArr = append(logArr, jobLog)
	}
	return
}
