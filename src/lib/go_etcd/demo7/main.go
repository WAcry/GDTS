package main

import (
	"context"
	"fmt"
	"lib/config"
	"time"
	"utils/coreos/etcd/clientv3"
	"utils/coreos/etcd/mvcc/mvccpb"
)

func main() {
	var (
		configs       clientv3.Config
		client        *clientv3.Client
		err           error
		kv            clientv3.KV
		watcher       clientv3.Watcher
		getResponse   *clientv3.GetResponse
		startRevision int64
		c             <-chan clientv3.WatchResponse
		watchResponse clientv3.WatchResponse
		event         *clientv3.Event
	)

	// set config
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

	// used to get, put, delete, watch key-value pairs
	kv = clientv3.NewKV(client)

	// simulate etcd KV changes
	go func() {
		for {
			kv.Put(context.TODO(), "/gdts/job7", "i am job7")

			kv.Delete(context.TODO(), "/gdts/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	// get current value
	if getResponse, err = kv.Get(context.TODO(), "/gdts/job7"); err != nil {
		fmt.Println(err)
		return
	}

	// if key-value exists here
	if len(getResponse.Kvs) != 0 {
		fmt.Println("current value:", string(getResponse.Kvs[0].Value))
	}

	// getResponse.Header.Revision: current etcd transaction ID, monotonic increasing
	startRevision = getResponse.Header.Revision + 1

	// create a watcher
	watcher = clientv3.NewWatcher(client)

	fmt.Println("watch event start from revision:", startRevision)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})

	// start watcher, canceled after 5 seconds
	c = watcher.Watch(ctx, "/gdts/job7", clientv3.WithRev(startRevision))

	// deal with event of key-value changes
	for watchResponse = range c {
		for _, event = range watchResponse.Events { // index is not used, so _ is ok
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("put:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("delete", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}
