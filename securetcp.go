package MagicCylinderG

import (
	"io"
	"net"
)

const bufSize = 1024

// Read 从连接中读取数据
func Read(conn *net.TCPConn, buf []byte) (int, error) {
	n, err := conn.Read(buf)
	return n, err
}

// Write 向连接中写入数据
func Write(conn *net.TCPConn, buf []byte) (int, error) {
	n, err := conn.Write(buf)
	return n, err
}

// Copy 不断的从src中读取数据写入到dst中
func Copy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, bufSize)
	for {
		readCount, err := Read(src, buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, err := Write(dst, buf[0:readCount])
			if err != nil {
				return err
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
