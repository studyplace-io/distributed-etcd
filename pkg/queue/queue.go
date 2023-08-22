package queue

import (
	"github.com/practice/etcd-distributed/pkg/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
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

func NewDistributeQueue(queueName string, path string) *DistributeQueue {
	dq := &DistributeQueue{
		etcdClient: client.EtcdClient(path),
		QueueName: queueName,
	}
	dq.Queue = recipe.NewQueue(dq.etcdClient, queueName)
	return dq
}
