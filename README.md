## distributed-etcd-sdk 
### 基于etcd封装分布式能力sdk
- 使用方法
    
    1. 修改配置文件config.yaml
       ```
       endpoint: http://127.0.0.1:2379
       prefix: /etcd-test
       ```
    2. 启动etcd实例
```bash
1. 分布式leader选主 
2. 分布式锁
3. 分布式读写锁
4. 分布式队列
5. 分布式优先队列
6. 分布式栅栏
```

### 使用方法

#### 分布式队列
```go
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

```

#### 分布式优先队列
```go
func TestPriorityQueue(t *testing.T) {
	dq := NewDistributePriorityQueue("my-priority-queue", "../../config.yaml")
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

```

#### 分布式锁
```go
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
```
#### 分布式栅栏
```go
func TestDistributeBarrier(t *testing.T) {

	bb := NewDistributeBarrier("my-barrier", "../../config.yaml")

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

```

#### 分布式leader选主 
```go
func TestLeaderElection(t *testing.T) {

	simulationProcess := func(myIdentity string) {
		le, err := NewLeaderElection(LeaderElectionConfig{
			EtcdClient: client.EtcdClient("../../config.yaml"),
			LeaseSeconds: 10,
			ElectionName: "leader-election-test",
			Identity: myIdentity,
			Callbacks: LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					log.Printf("OnStarted: %s is leader", myIdentity)
					time.Sleep(3 * time.Second)
					log.Printf("OnStarted: %s leader done", myIdentity)
				},
				OnStoppedLeading: func() {
					log.Printf("OnStopped: %s exit", myIdentity)
				},
				OnNewLeader: func(identity string) {
					if identity != myIdentity {
						log.Printf("OnNewLeader: leader from %s change to  %s", myIdentity, identity)
					} else {
						log.Printf("OnNewLeader: leader still is %s", identity)
					}
				},
			},
		})

		if err != nil {
			log.Fatalf("leader election error: %v", err)
		}
		defer le.Close()

		le.Run(context.Background())
	}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := fmt.Sprintf("simulation-process-%v", i)
			simulationProcess(id)
		}()
	}

	wg.Wait()



}

```