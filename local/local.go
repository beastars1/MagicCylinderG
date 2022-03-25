package local

import (
	"MagicCylinderG"
	"MagicCylinderG/cmd"
	"MagicCylinderG/crypto"
	"fmt"
	"log"
	"net"
	"strconv"
)

// local
// local指本机的监听地址，socks请求需要指定为该地址
// remote指服务端地址，socks请求加密后发送过去进行代理
type local struct {
	local  *net.TCPAddr
	remote *net.TCPAddr
	crypto crypto.Crypto
}

func NewLocal(localHost string, conf *cmd.Config) (*local, error) {
	// 接收conf，生成crypto
	c := crypto.CreateCrypto(conf)
	// 本机监听地址
	listenAddr := localHost + ":" + strconv.Itoa(conf.LocalPort)
	listener, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	remoteAddr := conf.RemoteHost + ":" + strconv.Itoa(conf.RemotePort)
	remote, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &local{
		local:  listener,
		remote: remote,
		crypto: c,
	}, nil
}

// Listen 监听本地的连接请求，didListen：连接建立之后的回调
func (l *local) Listen(didListen func(addr *net.TCPAddr)) error {
	return MagicCylinderG.ListenEncryptedConn(l.local, l.crypto, l.handleConn, didListen)
}

// 处理本机发起的socks请求，将其加密转发到服务端，并处理服务端的加密响应
func (l *local) handleConn(conn *MagicCylinderG.EncryptTcpConn) {
	defer conn.Close()
	// 连接到服务端
	remoteConn, err := MagicCylinderG.DialEncryptedConn(l.remote, l.crypto)
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
	conn.EncoderCopy(remoteConn)
}
