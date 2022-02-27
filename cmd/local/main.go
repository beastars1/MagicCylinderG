package main

import (
	"MagicCylinderG/cmd"
	"MagicCylinderG/local"
	"fmt"
	"log"
	"net"
)

const listenHost = "127.0.0.1"

func main() {
	conf := &cmd.Config{}
	// 读取配置
	conf.ReadConf()
	// 创建本机端
	spLocal, err := local.NewLocal(listenHost, conf.RemoteHost, conf.LocalPort, conf.RemotePort)
	if err != nil {
		log.Fatalln(err)
	}
	// 本机端监听请求
	log.Fatalln(spLocal.Listen(func(addr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
local 启动成功：
本地监听地址：
%s:%d
远程服务地址：
%s:%d
`, listenHost, conf.LocalPort, conf.RemoteHost, conf.RemotePort))
	}))
}
