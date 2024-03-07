package client

import "testing"

func TestEtcdClient(t *testing.T) {
	GetClientFromFileOrDie("../../config.yaml")
}
