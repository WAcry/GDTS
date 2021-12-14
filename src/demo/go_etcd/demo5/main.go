package main

import (
	"context"
	"demo/config"
	"fmt"
	"time"
	"utils/coreos/etcd/clientv3"
	"utils/coreos/etcd/mvcc/mvccpb"
)

func main() {
	var (
		conf    clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair  *mvccpb.KeyValue
	)

	conf = clientv3.Config{
		Endpoints: []string{config.ETCD_SERVER},
		//Endpoints:   []string{"21.281.122.24:2379", "21.281.122.39:2379", "21.281.122.21:2379"}, // cluster endpoints
		DialTimeout: 5 * time.Second,
	}

	// create connection
	if client, err = clientv3.New(conf); err != nil {
		fmt.Println(err)
		return
	}

	// used to read and write key-value pairs
	kv = clientv3.NewKV(client)

	// delete key-value pair
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV()); err != nil {
		// withFromKey() + withLimit(2) + key = /cron/jobs/job1 : delete /cron/jobs/job1, /cron/jobs/job2
		fmt.Println(err)
		return
	}

	println(1)

	// what is the deleted value?
	if len(delResp.PrevKvs) != 0 {
		for _, kvpair = range delResp.PrevKvs {
			fmt.Println("delete:", string(kvpair.Key), string(kvpair.Value))
		}
	}
}
