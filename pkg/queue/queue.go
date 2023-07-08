package queue

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"golanglearning/new_project/etcd-distributed/pkg/client"
)

// DistributeQueue 分布式队列
// 使用etcd中间件中转，可以支持多进程入队出队
type DistributeQueue struct {
	// etcdClient 客户端
	etcdClient  *clientv3.Client
	// 分布式队列
	*recipe.Queue
	// QueueName 名
	QueueName string
}

func NewDistributeQueue(queueName string) *DistributeQueue {
	dq := &DistributeQueue{
		etcdClient: client.EtcdClient(),
		QueueName: queueName,
	}
	dq.Queue = recipe.NewQueue(dq.etcdClient, queueName)
	return dq
}
