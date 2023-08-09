package priority_queue

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"golanglearning/new_project/etcd-distributed/pkg/client"
)

// DistributePriorityQueue 分布式优先队列
// 使用etcd中间件中转，可以支持多进程入队出队
type DistributePriorityQueue struct {
	// etcdClient 客户端
	etcdClient  *clientv3.Client
	// 分布式优先队列
	*recipe.PriorityQueue
	// QueueName 名
	QueueName string
}

func NewDistributePriorityQueue(queueName string, path string) *DistributePriorityQueue {
	dq := &DistributePriorityQueue{
		etcdClient: client.EtcdClient(path),
		QueueName: queueName,
	}
	dq.PriorityQueue = recipe.NewPriorityQueue(dq.etcdClient, queueName)
	return dq
}
