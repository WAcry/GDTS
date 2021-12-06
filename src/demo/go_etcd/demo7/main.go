package main

import (
	"config"
	"context"
	"fmt"
	"time"
	"utils/coreos/etcd/clientv3"
	"utils/coreos/etcd/mvcc/mvccpb"
)

func main() {
	var (
		conf               clientv3.Config
		client             *clientv3.Client
		err                error
		kv                 clientv3.KV
		watcher            clientv3.Watcher
		getResp            *clientv3.GetResponse
		watchStartRevision int64
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		event              *clientv3.Event
	)

	// set config
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

	// used to get, put, delete, watch key-value pairs
	kv = clientv3.NewKV(client)

	// simulate etcd KV changes
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")

			kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	// get current value
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// if key-value exists here
	if len(getResp.Kvs) != 0 {
		fmt.Println("current value:", string(getResp.Kvs[0].Value))
	}

	// getResp.Header.Revision: current etcd transaction ID, monotonic increasing
	watchStartRevision = getResp.Header.Revision + 1

	// create a watcher
	watcher = clientv3.NewWatcher(client)

	fmt.Println("watch event start from revision:", watchStartRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	// start watcher, canceled after 5 seconds
	watchRespChan = watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

	// deal with event of key-value changes
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events { // index is not used, so _ is ok
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("put:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("delete", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}
