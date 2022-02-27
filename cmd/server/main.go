package main

import (
	"MagicCylinderG/cmd"
	"MagicCylinderG/server"
	"fmt"
	"log"
	"net"
)

const listenHost = "127.0.0.1"

func main() {
	conf := &cmd.Config{}
	conf.ReadConf()
	s, err := server.NewServer(listenHost, conf.RemotePort)
	if err != nil {
		log.Fatalln(err)
	}
	s.Listen(func(addr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
server 启动成功：
本地监听地址：
%s:%d
`, listenHost, conf.LocalPort))
	})
}
