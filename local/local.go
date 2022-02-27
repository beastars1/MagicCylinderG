package local

import (
	"MagicCylinderG"
	"fmt"
	"log"
	"net"
	"strconv"
)

type SpLocal struct {
	listen *net.TCPAddr
	remote *net.TCPAddr
}

func NewLocal(listenHost, remoteHost string, localPort, remotePort int) (*SpLocal, error) {
	listenAddr := listenHost + ":" + strconv.Itoa(localPort)
	listen, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	remoteAddr := remoteHost + ":" + strconv.Itoa(remotePort)
	remote, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &SpLocal{
		listen: listen,
		remote: remote,
	}, nil
}

// Listen 监听本地的连接请求，didListen：连接建立之后的回调
func (l *SpLocal) Listen(didListen func(addr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", l.listen)
	if err != nil {
		return err
	}
	defer listener.Close()
	if didListen != nil {
		didListen(listener.Addr())
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
			continue
		}
		conn.SetLinger(0)
		go l.handleConn(conn)
	}
	return nil
}

func (l *SpLocal) handleConn(conn *net.TCPConn) {
	defer conn.Close()
	remoteConn, err := net.DialTCP("tcp", nil, l.remote)
	if err != nil {
		log.Fatal(fmt.Sprintf("连接到远程服务器 %s 失败:%s", l.remote, err))
		return
	}
	defer remoteConn.Close()
	remoteConn.SetLinger(0)

	// 转发,将服务端的响应转发到本地
	go func() {
		err := MagicCylinderG.Copy(conn, remoteConn)
		if err != nil {
			conn.Close()
			remoteConn.Close()
		}
	}()
	// 转发，将本机的请求转发到服务端
	MagicCylinderG.Copy(remoteConn, conn)
}
