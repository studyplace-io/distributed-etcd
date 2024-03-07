package locker

import (
	"github.com/practice/etcd-distributed/pkg/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
	"log"
	"math/rand"
	"sync"
	"time"
)

// DistributeRWLocker 分布式读写锁
// 使用etcd中间件中转，可以支持多进程争抢锁
type DistributeRWLocker struct {
	// etcdClient 客户端
	etcdClient *clientv3.Client
	// Locker 读写锁
	Locker *recipe.RWMutex
	// lockerName 名
	lockerName string
}

func NewDistributeRWLocker(lockerName string, path string) *DistributeRWLocker {
	l := &DistributeRWLocker{
		etcdClient: client.GetClientFromFileOrDie(path),
		lockerName: lockerName,
	}
	// 为锁生成session
	s1, err := concurrency.NewSession(l.etcdClient)
	if err != nil {
		log.Fatal(err)
	}
	l.Locker = recipe.NewRWMutex(s1, lockerName)

	return l
}

func doSomething1(id int, dl *DistributeRWLocker, wg *sync.WaitGroup) {
	defer wg.Done()

	// 请求锁
	log.Println("acquiring lock for ID:", id)
	dl.Locker.Lock()
	log.Println("acquired lock for ID:", id)

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	dl.Locker.Unlock()

	log.Println("released lock for ID:", id)
}

func writeSomething(id int, dl *DistributeRWLocker, wg *sync.WaitGroup) {
	defer wg.Done()

	// 请求锁
	log.Println("acquiring lock for ID:", id)
	if err := dl.Locker.Lock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired lock for ID:", id)

	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	if err := dl.Locker.Unlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released lock for ID:", id)
}

func readSomething(id int, dl *DistributeRWLocker, wg *sync.WaitGroup) {
	defer wg.Done()

	// 请求锁
	log.Println("acquiring rlock for ID:", id)
	if err := dl.Locker.RLock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired lock for ID:", id)

	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	if err := dl.Locker.RUnlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released rlock for ID:", id)
}
