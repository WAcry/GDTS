package main

import (
	"context"
	"demo/config"
	"fmt"
	"time"
	"utils/coreos/etcd/clientv3"
)

// implement a distributed optimistic lock
// lease is used to release the key-value lock when needed, which is safer and easier to implement
func main() {
	var (
		configs           clientv3.Config
		client            *clientv3.Client
		err               error
		lease             clientv3.Lease
		grantResponse     *clientv3.LeaseGrantResponse
		leaseId           clientv3.LeaseID
		c                 <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResponse *clientv3.LeaseKeepAliveResponse
		ctx               context.Context
		cancel            context.CancelFunc
		kv                clientv3.KV
		txn               clientv3.Txn
		txnResponse       *clientv3.TxnResponse
	)

	// client config
	configs = clientv3.Config{
		Endpoints: []string{config.ETCD_SERVER},
		//Endpoints:   []string{"21.281.122.24:2379", "21.281.122.39:2379", "21.281.122.21:2379"}, // cluster endpoints
		DialTimeout: 5 * time.Second,
	}

	// create connection
	if client, err = clientv3.New(configs); err != nil {
		fmt.Println(err)
		return
	}

	// use lease to implement lock auto expire
	// op operations
	// txn transaction: if else then

	// 1, lock (create a lease, auto renew, hold the lease to acquire a key)
	lease = clientv3.NewLease(client)

	// apply for a lease (5 seconds)
	if grantResponse, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	}

	// get lease id
	leaseId = grantResponse.ID

	// create a context used to cancel lease keep alive
	ctx, cancel = context.WithCancel(context.TODO())

	// make sure to cancel lease keep alive after return and exit
	defer cancel()
	defer lease.Revoke(context.TODO(), leaseId)

	// cancel lease keep alive after a while 5 seconds
	if c, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// deal with keep alive response
	go func() {
		for {
			select {
			case keepAliveResponse = <-c:
				if c == nil {
					fmt.Println("lease expired")
					goto END
				} else {
					fmt.Println("keep alive once, lease id:", keepAliveResponse.ID)
				}
			}
		}
	END:
	}()

	// if key not exists, set it, else lock failed
	kv = clientv3.NewKV(client)

	// create transaction
	txn = kv.Txn(context.TODO())

	// define transaction

	// if key not exists
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock9")) // else fail to lock

	// submit transaction
	if txnResponse, err = txn.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	// if failed to acquire lock
	if !txnResponse.Succeeded {
		fmt.Println("lock is acquired by others:", string(txnResponse.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2, do the job

	fmt.Println("do the job")
	time.Sleep(5 * time.Second)

	// 3, unlock (cancel auto renew, release lease)
	// this is done by defer above
	// they will release the lease and cancel the keep alive, delete related key-value
}
