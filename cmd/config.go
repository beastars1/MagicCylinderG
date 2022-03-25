package cmd

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	LocalPort  int    `json:"localPort"`
	RemoteHost string `json:"remoteHost"`
	RemotePort int    `json:"remotePort"`
}

func (conf *Config) ReadConf() {
	configPath := "config.json"
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("open %s error:%s", configPath, err)
		}
		defer file.Close()
		err = json.NewDecoder(file).Decode(conf)
		if err != nil {
			log.Fatalf("json conf error:%s", err)
		}
	}
}

// DefaultConf 如果没有配置文件，生成默认配置
func DefaultConf() *Config {
	return &Config{
		LocalPort:  9999,
		RemoteHost: "127.0.0.1",
		RemotePort: 8888,
	}
}
