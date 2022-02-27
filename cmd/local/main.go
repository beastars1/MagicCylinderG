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
	conf.ReadConf()
	spLocal, err := local.NewLocal(listenHost, conf.RemoteHost, conf.LocalPort, conf.RemotePort)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(spLocal.Listen(func(addr net.Addr) {
		log.Println(fmt.Sprintf(`
local 启动成功：
本地监听地址：
%s:%d
远程服务地址：
%s:%d
`, listenHost, conf.LocalPort, conf.RemoteHost, conf.RemotePort))
	}))
}
