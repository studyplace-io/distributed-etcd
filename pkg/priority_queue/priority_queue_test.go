package priority_queue

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestPriorityQueue(t *testing.T) {
	dq := NewDistributePriorityQueue("my-priority-queue")
	defer dq.etcdClient.Close()

	wg := &sync.WaitGroup{}

	// 模拟入队
	for i := 0; i < 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			dq.Enqueue(fmt.Sprintf("key-%d", i), uint16(i*100+i))
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
