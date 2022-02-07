package worker

import (
	"common"
	"context"
	"time"
	"utils/coreos/etcd/clientv3"
	"utils/coreos/etcd/mvcc/mvccpb"
)

type WkJobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	JobManager *WkJobManager
)

// listen job changes
func (jobMgr *WkJobManager) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	// get all jobs under /gdts/jobs/, and get the current revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// current jobs
	for _, kvpair = range getResp.Kvs {
		// decode job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// sync jobEvent to scheduler
			Scheduler.PushJobEvent(jobEvent)
		}
	}

	// 2, start listen from this revision
	go func() {
		watchStartRevision = getResp.Header.Revision + 1
		// listen /gdts/jobs/ dir
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // save task event
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE: // delete task event
					// Delete /gdts/jobs/job10
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				// push jobEvent to scheduler
				Scheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// listen killer
func (jobMgr *WkJobManager) watchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
		job        *common.Job
	)
	// listen /gdts/killer dir
	go func() {
		// listen /gdts/killer/ dir
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // kill task event
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)
					// push jobEvent to scheduler
					Scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: // killer is expired, auto deleted
				}
			}
		}
	}()
}

func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	config = clientv3.Config{
		Endpoints:   Config.EtcdEndpoints,
		DialTimeout: time.Duration(Config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	JobManager = &WkJobManager{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	JobManager.watchJobs()

	JobManager.watchKiller()

	return
}

func (jobMgr *WkJobManager) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
