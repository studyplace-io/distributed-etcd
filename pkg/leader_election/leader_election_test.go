package leader_election

import (
	"context"
	"fmt"
	"github.com/practice/etcd-distributed/pkg/client"
	"log"
	"sync"
	"testing"
	"time"
)

func TestLeaderElection(t *testing.T) {

	simulationProcess := func(myIdentity string) {
		le, err := NewLeaderElection(LeaderElectionConfig{
			EtcdClient:   client.GetClientFromFileOrDie("../../config.yaml"),
			LeaseSeconds: 10,
			ElectionName: "leader-election-test",
			Identity:     myIdentity,
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
