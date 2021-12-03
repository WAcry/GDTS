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
		conf    clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
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

	// used to get, put, delete, watch key-value pairs
	kv = clientv3.NewKV(client)

	// write another job
	kv.Put(context.TODO(), "/cron/jobs/job2", "{...}")

	// read all keys with prefix /cron/jobs/
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else { // successfully get all keys
		fmt.Println(getResp.Kvs, getResp.Count)
	}
}
