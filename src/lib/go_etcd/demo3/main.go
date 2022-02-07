package main

import (
	"context"
	"fmt"
	"lib/config"
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

	if response, err = kv.Get(context.TODO(), "/gdts/job1"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.Kvs, response.Count)
	}
}
