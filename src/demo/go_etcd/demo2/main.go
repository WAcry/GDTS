package main

import (
	"context"
	"demo/config"
	"fmt"
	"time"
	"utils/coreos/etcd/clientv3"
)

func main() {
	var (
		configs  clientv3.Config
		client   *clientv3.Client
		err      error
		kv       clientv3.KV
		response *clientv3.PutResponse
	)

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

	// used to read and write key-value pairs
	kv = clientv3.NewKV(client)

	if response, err = kv.Put(context.TODO(), "/cron/job1", "bye", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", response.Header.Revision)
		if response.PrevKv != nil { // print nothing first (no previous value), then print bye next execution
			fmt.Println("PrevValue:", string(response.PrevKv.Value))
		}
	}
}
