package worker

import (
	"common"
	"context"
	"net"
	"time"
	"utils/coreos/etcd/clientv3"
)

// WkRegister register the server to etcd: /gdts/workers/<IP>
type WkRegister struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string
}

var (
	Register *WkRegister
)

func getLocalIP() (ipv4 string, err error) {
	var (
		addrList []net.Addr
		addr     net.Addr
		ipNet    *net.IPNet // IP address
		isIpNet  bool
	)
	// get all network interfaces
	if addrList, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// get first non-loopback interface
	for _, addr = range addrList {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// skip ipv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}

	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

// register to /gdts/workers/<IP>
func (register *WkRegister) keepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)

	for {
		regKey = common.JOB_WORKER_DIR + register.localIP

		cancelFunc = nil

		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}

		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { // failed to keep alive
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}

func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)

	config = clientv3.Config{
		Endpoints:   Config.EtcdEndpoints,
		DialTimeout: time.Duration(Config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	if localIp, err = getLocalIP(); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	Register = &WkRegister{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}

	// register the server to etcd
	go Register.keepOnline()

	return
}
