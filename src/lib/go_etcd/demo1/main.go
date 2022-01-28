package main

import (
	"fmt"
	"lib/config"
	"time"
	"utils/coreos/etcd/clientv3"
)

func main() {
	var (
		configs clientv3.Config
		client  *clientv3.Client
		err     error
	)

	// client config
	configs = clientv3.Config{
		Endpoints: []string{config.ETCD_SERVER},
		//Endpoints:   []string{"21.281.122.24:2379", "21.281.122.39:2379", "21.281.122.21:2379"}, // cluster endpoints
		DialTimeout: 10 * time.Second,
	}

	// create connection
	if client, err = clientv3.New(configs); err != nil {
		fmt.Println(err)
		return
	}

	client = client
}
