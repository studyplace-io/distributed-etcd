package client

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

func GetClientFromFileOrDie(path string) *clientv3.Client {
	config, err := loadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return GetClientOrDie(config)
}

func GetClientOrDie(config *EtcdConfig) *clientv3.Client {

	endpoints := []string{config.Endpoint}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}

	return cli
}
