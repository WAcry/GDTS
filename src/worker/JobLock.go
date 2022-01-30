package worker

import (
	"common"
	"context"
	"utils/coreos/etcd/clientv3"
)

// JobLock distributed lock (TXN transaction)
type JobLock struct {
	kv    clientv3.KV
	lease clientv3.Lease

	jobName    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	isLocked   bool
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)

	// 1, create lease for 5s
	if leaseGrantResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	leaseId = leaseGrantResp.ID

	// 2, keep lease alive
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	// 3, deal with lease keepalive response
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	// 4, create txn transaction
	txn = jobLock.kv.Txn(context.TODO())

	lockKey = common.JOB_LOCK_DIR + jobLock.jobName

	// 5, acquire lock in txn transaction
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// submit txn transaction
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	// 6, return if success, revoke lease if failed
	if !txnResp.Succeeded {
		err = common.ERR_LOCK_ALREADY_ACQUIRED
		goto FAIL
	}

	// success acquire lock
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true
	return

FAIL:
	cancelFunc()                                  // cancel lease keepalive
	jobLock.lease.Revoke(context.TODO(), leaseId) //  release lease
	return
}

func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
		jobLock.cancelFunc()                                  // cancel lease keepalive
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseId) // release lease
	}
}
