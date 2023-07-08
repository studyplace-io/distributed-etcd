package client

import (
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"
)


func EtcdClient() *clientv3.Client {
	config, err := loadConfig("../../config.yaml")
	endpoints := []string{config.Endpoint}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
