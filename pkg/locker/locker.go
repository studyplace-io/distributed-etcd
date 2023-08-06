package locker

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golanglearning/new_project/etcd-distributed/pkg/client"
	"log"
	"math/rand"
	"sync"
	"time"
)

// DistributeLocker 分布式锁
// 使用etcd中间件中转，可以支持多进程争抢锁
type DistributeLocker struct {
	// etcdClient 客户端
	etcdClient  *clientv3.Client
	// Locker 锁
	Locker sync.Locker
	// lockerName 名
	lockerName string
}

func NewDistributeLocker(lockerName string) *DistributeLocker {
	l := &DistributeLocker{
		etcdClient: client.EtcdClient("../../config.yaml"),
		lockerName: lockerName,
	}
	// 为锁生成session
	s1, err := concurrency.NewSession(l.etcdClient)
	if err != nil {
		log.Fatal(err)
	}
	// l.Locker := concurrency.NewMutex(s1, lockName)
	l.Locker = concurrency.NewLocker(s1, lockerName)

	return l
}


func doSomething(id int, dl *DistributeLocker, wg *sync.WaitGroup) {
	defer wg.Done()

	// 请求锁
	log.Println("acquiring lock for ID:", id)
	dl.Locker.Lock()
	log.Println("acquired lock for ID:", id)

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	dl.Locker.Unlock()

	log.Println("released lock for ID:", id)
}
