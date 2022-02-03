package worker

import (
	"common"
	"context"
	"time"
	"utils/mongodb/mongo-go-driver/mongo"
	"utils/mongodb/mongo-go-driver/mongo/clientopt"
)

// WkLogSink mongodb log storer
type WkLogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	// LogSink singleton logger
	LogSink *WkLogSink
)

// write logs in batches
func (logSink *WkLogSink) saveLogs(batch *common.LogBatch) {
	logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
}

func (logSink *WkLogSink) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch // current batch
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch // timeout batch
	)

	for {
		select {
		case log = <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// let this batch of logs auto commit after 1 second
				commitTimer = time.AfterFunc(
					time.Duration(Config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}

			// append new log to current batch
			logBatch.Logs = append(logBatch.Logs, log)

			// if batch is full, send automatically
			if len(logBatch.Logs) >= Config.JobLogBatchSize {
				// send log batch to mongodb
				logSink.saveLogs(logBatch)
				// clean current batch
				logBatch = nil
				// cancel timer
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan: // timeout batch
			// check if timeout batch is current batch
			if timeoutBatch != logBatch {
				continue // skip if submitted
			}
			// write timeout batch to mongodb
			logSink.saveLogs(timeoutBatch)
			// clear current batch
			logBatch = nil // TODO
		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)

	if client, err = mongo.Connect(
		context.TODO(),
		Config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(Config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		return
	}

	LogSink = &WkLogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go LogSink.writeLoop()
	return
}

// Append send log
func (logSink *WkLogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// if queue is full, drop log
	}
}
