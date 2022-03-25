package crypto

import (
	"encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Simple struct {
	encodePassword *password
	decodePassword *password
}

func (crypto *Simple) Decode(data []byte) {
	for i, v := range data {
		data[i] = crypto.decodePassword[v]
	}
}

func (crypto *Simple) Encode(data []byte) {
	for i, v := range data {
		data[i] = crypto.encodePassword[v]
	}
}

func NewSimpleCrypto(auth string) *Simple {
	encodePassword, _ := parse2Byte(auth)
	decodePassword := &password{}
	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	return &Simple{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}

const pwdLen = 256

type password [pwdLen]byte

func init() {
	rand.Seed(time.Now().Unix())
}

// 将byte的密码使用base64加密为字符串
func (p *password) String() string {
	s := base64.StdEncoding.EncodeToString(p[0:pwdLen])
	return s
}

// RandPassword 随机生成一串可使用的密码字符串
func randPassword() {
	// 生成一个0~255的不重复随机数组
	bytes := rand.Perm(pwdLen)
	pwd := &password{}
	for i, v := range bytes {
		pwd[i] = byte(v)
		if i == v {
			randPassword()
			return
		}
	}
	log.Fatalf("simple auth:%s", pwd.String())
}

// Parse2Byte 将给定的字符串base64解码为byte密码
func parse2Byte(pwdStr string) (*password, error) {
	bytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(pwdStr))
	if err != nil || len(bytes) != pwdLen {
		randPassword()
	}
	pwd := password{}
	copy(pwd[:], bytes)
	bytes = nil
	return &pwd, nil
}
