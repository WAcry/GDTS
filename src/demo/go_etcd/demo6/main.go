package main

import (
	"config"
	"context"
	"fmt"
	"time"
	"utils/coreos/etcd/clientv3"
)

func main() {
	var (
		conf           clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		kv             clientv3.KV
	)

	conf = clientv3.Config{
		Endpoints: []string{config.Server["etcd"]},
		//Endpoints:   []string{"21.281.122.24:2379", "21.281.122.39:2379", "21.281.122.21:2379"}, // cluster endpoints
		DialTimeout: 5 * time.Second,
	}

	// create connection
	if client, err = clientv3.New(conf); err != nil {
		fmt.Println(err)
		return
	}

	// acquire lease
	lease = clientv3.NewLease(client)

	// apply a lease of 10 seconds
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	// get lease id
	leaseId = leaseGrantResp.ID

	// set auto keepalive to lease, interval is 1/3 of lease time, so (10/3) ≈ 3 seconds
	// so the lease will be automatically renewed and not expire, until context is cancelled after 20 seconds
	// so keepalive will be cancelled after 20 seconds, and the lease will be expired after about 10+20=30 seconds
	ctx, _ := context.WithTimeout(context.TODO(), 20*time.Second)
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
		return
	}

	// deal with every response of keepalive request
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					fmt.Println("keepAlive is canceled", " at: ", time.Now())
					goto END
				} else { // request keepalive and get one response every second
					fmt.Println("lease receives keepalive request: ", keepResp.ID, " at: ", time.Now())
				}
			}
		}
	END:
	}()

	// get kv API
	kv = clientv3.NewKV(client)

	// Put a key-value pair with lease, so that it expires after 10 seconds
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("write successfully:", putResp.Header.Revision)

	// check if the key is expired every 2 seconds
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv is expired", " at: ", time.Now())
			break
		}
		fmt.Println("not expired yet:", getResp.Kvs, " at: ", time.Now())
		time.Sleep(2 * time.Second)
	}
}
