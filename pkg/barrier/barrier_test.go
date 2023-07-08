package barrier

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestDistributeBarrier(t *testing.T) {

	bb := NewDistributeBarrier("my-barrier")

	err := bb.Hold()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		i := i
		go func() {

			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			log.Println("enter for ID:", i)
			err := bb.Wait()
			if err != nil {
				panic(err)
			}
			log.Println("entered for ID:", i)
			wg.Done()
		}()
	}

	time.Sleep(12 * time.Second)
	err = bb.Release()
	if err != nil {
		panic(err)
	}

	wg.Wait()
}


