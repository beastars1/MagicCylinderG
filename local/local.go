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
func (l *SpLocal) Listen(didListen func(addr *net.TCPAddr)) error {
	return MagicCylinderG.ListenEncryptedConn(l.listen, l.handleConn, didListen)
}

func (l *SpLocal) handleConn(conn *MagicCylinderG.EncryptTcpConn) {
	defer conn.Close()
	// 连接到服务端
	remoteConn, err := MagicCylinderG.DialEncryptedConn(l.remote)
	if err != nil {
		log.Println(fmt.Sprintf("连接到远程服务器 %s 失败:%s", l.remote, err))
		return
	}
	defer remoteConn.Close()

	// 转发,将服务端的响应解密后转发到本地
	go func() {
		err := remoteConn.DecoderCopy(conn)
		if err != nil {
			conn.Close()
			remoteConn.Close()
		}
	}()
	// 转发，将本机的请求加密转发到服务端
	remoteConn.EncoderCopy(conn)
}
