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
		response *clientv3.GetResponse
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

	// used to get, put, delete, watch key-value pairs
	kv = clientv3.NewKV(client)

	// write another job
	kv.Put(context.TODO(), "/cron/job2", "{...}")

	// read all keys with prefix /cron/
	if response, err = kv.Get(context.TODO(), "/cron/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else { // successfully get all keys
		fmt.Println(response.Kvs, response.Count)
	}
}
