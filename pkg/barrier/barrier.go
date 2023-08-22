package barrier

import (
	"github.com/practice/etcd-distributed/pkg/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

// DistributeBarrier 分布式队列
// 使用etcd中间件中转，可以支持多进程入队出队
type DistributeBarrier struct {
	// etcdClient 客户端
	etcdClient  *clientv3.Client
	*recipe.Barrier
	// BarrierName 名
	BarrierName string
}

func NewDistributeBarrier(barrierName string, path string) *DistributeBarrier {
	dq := &DistributeBarrier{
		etcdClient: client.EtcdClient(path),
		BarrierName: barrierName,
	}
	dq.Barrier = recipe.NewBarrier(dq.etcdClient, barrierName)
	return dq
}
