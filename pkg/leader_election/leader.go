package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3/concurrency"
	"golanglearning/new_project/etcd-distributed/pkg/client"
	"log"
	"os"
)

var (
	nodeID    = flag.Int("id", 0, "node ID")
	electName = flag.String("name", "my-test-elect", "election name")
)

func main() {
	flag.Parse()

	cli := client.EtcdClient()

	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}

	e1 := concurrency.NewElection(s1, *electName)

	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "elect":
			go elect(e1, *electName)
		case "proclaim":
			proclaim(e1, *electName)
		case "resign":
			resign(e1, *electName)
		case "watch":
			go watch(e1, *electName)
		case "query":
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}

var count int

func elect(e1 *concurrency.Election, electName string) {
	log.Println("acampaigning for ID:", *nodeID)
	if err := e1.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("campaigned for ID:", *nodeID)
	count++
}

func proclaim(e1 *concurrency.Election, electName string) {
	log.Println("proclaiming for ID:", *nodeID)
	if err := e1.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("proclaimed for ID:", *nodeID)
	count++
}

func resign(e1 *concurrency.Election, electName string) {
	log.Println("resigning for ID:", *nodeID)
	if err := e1.Resign(context.TODO()); err != nil {
		log.Println(err)
	}
	log.Println("resigned for ID:", *nodeID)
}

func watch(e1 *concurrency.Election, electName string) {
	ch := e1.Observe(context.TODO())

	log.Println("start to watch for ID:", *nodeID)
	for i := 0; i < 10; i++ {
		resp := <-ch
		log.Println("leader changed to", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
	}
}

func query(e1 *concurrency.Election, electName string) {
	resp, err := e1.Leader(context.Background())
	if err != nil {
		log.Printf("failed to get the current leader: %v", err)
	}
	log.Println("current leader:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
}

func rev(e1 *concurrency.Election, electName string) {
	rev := e1.Rev()
	log.Println("current rev:", rev)
}