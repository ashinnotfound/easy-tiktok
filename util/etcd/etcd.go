package etcd

//etcd 是一个分布式键值对存储，设计用来可靠而快速的保存关键数据并提供访问。
//通过分布式锁，leader选举和写屏障(write barriers)来实现可靠的分布式协作。
//etcd集群是为高可用，持久性数据存储和检索而准备。
//使用
import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"log"
	"time"
)

var client *clientv3.Client

const (
	prefix = "service"
)

func init() {
	var err error
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.21.23.42:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
}

func Register(ctx context.Context, serviceName, addr string) error {
	log.Println("Try register to etcd ...")
	// 创建一个租约
	lease := clientv3.NewLease(client)
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	leaseResp, err := lease.Grant(cancelCtx, 3)
	if err != nil {
		return err
	}

	leaseChannel, err := lease.KeepAlive(ctx, leaseResp.ID) // 长链接, 不用设置超时时间
	if err != nil {
		return err
	}

	em, err := endpoints.NewManager(client, prefix)
	if err != nil {
		return err
	}

	cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	if err := em.AddEndpoint(cancelCtx, fmt.Sprintf("%s/%s/%s", prefix, serviceName, uuid.NewV4().String()), endpoints.Endpoint{
		Addr: addr,
	}, clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}
	log.Println("Register etcd success")

	del := func() {
		log.Println("Register close")

		cancelCtx, cancel = context.WithTimeout(ctx, time.Second*3)
		defer cancel()
		em.DeleteEndpoint(cancelCtx, serviceName)

		lease.Close()
	}
	// 保持注册状态(连接断开重连)
	keepRegister(ctx, leaseChannel, del, serviceName, addr)

	return nil
}

func keepRegister(ctx context.Context, leaseChannel <-chan *clientv3.LeaseKeepAliveResponse, cleanFunc func(), serviceName, addr string) {
	go func() {
		failedCount := 0
		for {
			select {
			case resp := <-leaseChannel:
				if resp != nil {
					//log.Println("keep alive success.")
				} else {
					log.Println("keep alive failed.")
					failedCount++
					for failedCount > 3 {
						cleanFunc()
						if err := Register(ctx, serviceName, addr); err != nil {
							time.Sleep(time.Second)
							continue
						}
						return
					}
					continue
				}
			case <-ctx.Done():
				cleanFunc()
				client.Close()
				return
			}
		}
	}()
}
