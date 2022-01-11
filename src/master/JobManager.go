package master

import (
	"common"
	"context"
	"encoding/json"
	"time"
	"utils/coreos/etcd/clientv3"
	"utils/coreos/etcd/mvcc/mvccpb"
)

type MstJobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	// JobManager singleton jobManager
	JobManager *MstJobManager
)

func InitJobManager() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	config = clientv3.Config{
		Endpoints:   Config.EtcdEndpoints,
		DialTimeout: time.Duration(Config.EtcdDialTimeout) * time.Millisecond,
	}

	// create client and connection
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// get API subsets of KV and Lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	JobManager = &MstJobManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

func (jobMgr *MstJobManager) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// save job into /cron/jobs/<jobname> as json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)

	// key of etcd
	jobKey = common.JOB_SAVE_DIR + job.Name
	// job info json
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	// save into etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	// if update, return old job
	if putResp.PrevKv != nil {
		// decode old job
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (jobMgr *MstJobManager) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj common.Job
	)

	// key of etcd
	jobKey = common.JOB_SAVE_DIR + name

	// delete from etcd
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// return deleted job
	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

func (jobMgr *MstJobManager) ListJobs() (jobList []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	dirKey = common.JOB_SAVE_DIR

	// get all jobs under the dir
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	jobList = make([]*common.Job, 0)

	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

func (jobMgr *MstJobManager) KillJob(name string) (err error) {
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	// workers listen on /cron/workers/ and will kill the <jobname>
	killerKey = common.JOB_KILLER_DIR + name

	// create a lease outdated after 1 second for <jobname>
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	leaseId = leaseGrantResp.ID

	// put <jobname> into /cron/killer/<jobname>
	if _, err = jobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
