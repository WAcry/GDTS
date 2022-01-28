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
		putOp    clientv3.Op
		getOp    clientv3.Op
		response clientv3.OpResponse
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

	kv = clientv3.NewKV(client)

	// create Op: operation
	putOp = clientv3.OpPut("/cron/job8", "123456789")

	// execute OP
	if response, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}

	// kv.Do(op)
	// ==
	// kv.Put
	// kv.Get
	// kv.Delete

	fmt.Println("write Revision:", response.Put().Header.Revision)

	// create Op
	getOp = clientv3.OpGet("/cron/job8")

	// execute OP
	if response, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}

	// print
	fmt.Println("create Revision:", response.Get().Kvs[0].CreateRevision)    // <= writeRevision
	fmt.Println("current data Revision:", response.Get().Kvs[0].ModRevision) // equal to write Revision
	fmt.Println("current data value:", string(response.Get().Kvs[0].Value))
}
