package client

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)


func EtcdClient(path string) *clientv3.Client {
	config, err := loadConfig(path)
	endpoints := []string{config.Endpoint}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
