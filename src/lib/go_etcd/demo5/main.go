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
		configs  clientv3.Config
		client   *clientv3.Client
		err      error
		kv       clientv3.KV
		response *clientv3.DeleteResponse
		pair     *mvccpb.KeyValue
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

	// delete key-value pair
	if response, err = kv.Delete(context.TODO(), "/cron/job1", clientv3.WithPrevKV()); err != nil {
		// withFromKey() + withLimit(2) + key = /cron/job1 : delete /cron/job1, /cron/job2
		fmt.Println(err)
		return
	}

	println(1)

	// what is the deleted value?
	if len(response.PrevKvs) != 0 {
		for _, pair = range response.PrevKvs {
			fmt.Println("delete:", string(pair.Key), string(pair.Value))
		}
	}
}
