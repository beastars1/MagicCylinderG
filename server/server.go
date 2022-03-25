package server

import (
	"MagicCylinderG"
	"MagicCylinderG/cmd"
	"MagicCylinderG/crypto"
	"encoding/binary"
	"net"
	"strconv"
)

// server local指服务端监听客户端请求的地址，默认127.0.0.1
type server struct {
	local  *net.TCPAddr
	crypto crypto.Crypto
}

func NewServer(localHost string, conf *cmd.Config) (*server, error) {
	// 接收conf，生成crypto
	c := crypto.CreateCrypto(conf)
	// 服务端监听地址
	listenAddr := localHost + ":" + strconv.Itoa(conf.RemotePort)
	listener, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return &server{
		local:  listener,
		crypto: c,
	}, err
}

// Listen  对监听的请求进行加密解密处理
func (s *server) Listen(didListen func(addr *net.TCPAddr)) error {
	return MagicCylinderG.ListenEncryptedConn(s.local, s.crypto, s.handleEncryptedConn, didListen)
}

// 处理服务端接收到的客户端的加密socks请求，进行socks握手，和客户端进行socks连接，之连接成功之后客户端就通过socks发送请求
// socks中传输的数据是加密的
func (s *server) handleEncryptedConn(conn *MagicCylinderG.EncryptTcpConn) {
	defer conn.Close()
	buf := make([]byte, 256)
	// 解析socks协议
	// 客户端发起握手请求
	_, err := conn.DecoderRead(buf)
	/*
	  |VER | NMETHODS | METHODS  |
	  | 1  |    1     | 1 to 255 |
	*/
	if err != nil || buf[0] != 0x05 {
		// 如果不是socks5协议，不建立连接
		return
	}
	/**
	  |VER | METHOD |
	  | 1  |   1    |
	*/
	// 响应socks5握手
	conn.EncoderWrite([]byte{0x05, 0x00})

	// 握手成功之后，客户端发送目标服务器的信息，由代理服务器尝试连接
	/**
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	*/
	n, err := conn.DecoderRead(buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	if buf[1] != 0x01 {
		return
	}
	atyp := buf[3]
	var ip []byte
	switch atyp {
	case 0x01:
		// ipv4
		ip = buf[4 : 4+net.IPv4len]
	case 0x03:
		// domain name
		// string(buf[4])是 \r，往后才是域名
		addr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			return
		}
		ip = addr.IP
	case 0x04:
		// ipv6
		ip = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	port := buf[n-2 : n]
	// 目标主机的地址
	addr := &net.TCPAddr{
		IP:   ip,
		Port: int(binary.BigEndian.Uint16(port)),
	}
	// 代理服务器连接目标主机
	dstConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	} else {
		defer dstConn.Close()
		dstConn.SetLinger(0)
		// 向客户端响应连接成功，可以正式通过socks发送请求了
		/**
		  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		  | 1  |  1  | X'00' |  1   | Variable |    2     |
		*/
		conn.EncoderWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	// socks连接建立成功，对请求进行转发
	// 从客户端发送的请求解密发送到目标主机
	go func() {
		err := conn.DecoderCopy(dstConn)
		if err != nil {
			// 转发过程中出现错误，直接关闭
			conn.Close()
			dstConn.Close()
		}
	}()

	// 从目标服务器响应的数据加密发送到客户端
	(&MagicCylinderG.EncryptTcpConn{
		ReadWriteCloser: dstConn,
		Crypto:          conn.Crypto,
	}).EncoderCopy(conn)
}
