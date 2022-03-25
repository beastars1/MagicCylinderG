package MagicCylinderG

import (
	"MagicCylinderG/crypto"
	"io"
	"log"
	"net"
	"sync"
)

const bufSize = 1024

var bpool sync.Pool

func init() {
	bpool.New = func() interface{} {
		return make([]byte, bufSize)
	}
}

func bufPoolGet() []byte {
	return bpool.Get().([]byte)
}

func bufPoolPut(buf []byte) {
	bpool.Put(buf)
}

type EncryptTcpConn struct {
	io.ReadWriteCloser
	crypto.Crypto
}

// DecoderRead 将conn中的数据解密后写入到buf
func (conn *EncryptTcpConn) DecoderRead(buf []byte) (int, error) {
	n, err := conn.Read(buf)
	// decoder buf
	conn.Decode(buf[:n])
	return n, err
}

// EncoderWrite 将buf加密后写入到conn
func (conn *EncryptTcpConn) EncoderWrite(buf []byte) (int, error) {
	// encoder buf
	conn.Encode(buf)
	n, err := conn.Write(buf)
	return n, err
}

// DecoderCopy 不断的从conn中读取数据解码写入到dst中
func (conn *EncryptTcpConn) DecoderCopy(dst io.ReadWriteCloser) error {
	buf := bufPoolGet()
	defer bufPoolPut(buf)
	for {
		// 从conn读取数据并解密到buf
		readCount, err := conn.DecoderRead(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		if readCount > 0 {
			// 将解密后的数据写入到conn
			writeCount, err := dst.Write(buf[0:readCount])
			if err != nil {
				return err
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// EncoderCopy 不断的从conn中读取数据加密写入到dst中
func (conn *EncryptTcpConn) EncoderCopy(dst io.ReadWriteCloser) error {
	buf := bufPoolGet()
	defer bufPoolPut(buf)
	for {
		// 从conn读取请求并写入到到buf
		readCount, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		if readCount > 0 {
			// 将数据加密写入到conn
			writeCount, err := (&EncryptTcpConn{
				ReadWriteCloser: dst,
				Crypto:          conn.Crypto,
			}).EncoderWrite(buf[0:readCount])
			if err != nil {
				return err
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// ListenEncryptedConn 监听连接到addr的请求，通过handleConn进行处理。
// 对于local端来说，local就是本机发起的请求，远端是服务器；
// 对于server端来说，local就是local端发送的加密请求，远端就是要访问的网站。
func ListenEncryptedConn(localAddr *net.TCPAddr, crypto crypto.Crypto, handleConn func(local *EncryptTcpConn), didListen func(listenAddr *net.TCPAddr)) error {
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	if didListen != nil {
		go didListen(listener.Addr().(*net.TCPAddr))
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		localConn.SetLinger(0)
		go handleConn(&EncryptTcpConn{
			ReadWriteCloser: localConn,
			Crypto:          crypto,
		})
	}
}

// DialEncryptedConn 连接到远程服务器
func DialEncryptedConn(remoteAddr *net.TCPAddr, crypto crypto.Crypto) (*EncryptTcpConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		return nil, err
	}
	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	remoteConn.SetLinger(0)
	return &EncryptTcpConn{
		ReadWriteCloser: remoteConn,
		Crypto:          crypto,
	}, err
}
