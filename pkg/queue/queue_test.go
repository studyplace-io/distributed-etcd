package queue

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	dq := NewDistributeQueue("my-queue", "../../config.yaml")
	defer dq.etcdClient.Close()

	wg := &sync.WaitGroup{}

	// 模拟入队
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			dq.Enqueue(fmt.Sprintf("key-%d", i))
		}()
	}

	// 模拟出队
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			v, err := dq.Dequeue()
			time.Sleep(time.Second)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("received: %s\n", v)
		}()
	}

	wg.Wait()

}
