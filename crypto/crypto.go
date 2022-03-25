package crypto

import "MagicCylinderG/cmd"

type Crypto interface {
	Decode(data []byte)
	Encode(data []byte)
}

func CreateCrypto(conf *cmd.Config) Crypto {
	switch conf.Method {
	case "none":
		return NewNoneCrypto()
	case "simple":
		return NewSimpleCrypto(conf.Auth)
	default:
		return NewNoneCrypto()
	}
}
