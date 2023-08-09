package locker

import (
	"sync"
	"testing"
)

func TestLocker(t *testing.T) {

	l := NewDistributeLocker("my-locker", "../../config.yaml")
	defer l.etcdClient.Close()

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go doSomething(i, l, &wg)
	}

	wg.Wait()
}

func TestRWLocker(t *testing.T) {
	rwl := NewDistributeRWLocker("my-rwlocker", "../../config.yaml")
	defer rwl.etcdClient.Close()

	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go readSomething(i, rwl, &wg)
	}

	for i := 0; i < 5; i++ {
		go writeSomething(10+i, rwl, &wg)
	}

	wg.Wait()
}
