package client

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"k8s.io/klog"
)

type EtcdConfig struct {
	Endpoint  string  `json:"endpoint"`
	Prefix    string  `json:"prefix"`
}

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{}
}

func loadConfigFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		klog.Errorf("load file err: %s", err)
		return nil
	}
	return b
}

// LoadConfig 读取配置文件
func loadConfig(path string) (*EtcdConfig, error) {
	c := NewEtcdConfig()
	if b := loadConfigFile(path); b != nil {
		err := yaml.Unmarshal(b, c)
		if err != nil {
			klog.Errorf("unmarshal err: %s", err)
			return nil, err
		}
		return c, err
	} else {
		return nil, fmt.Errorf("load config file error")
	}
}

