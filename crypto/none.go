package crypto

// None 不进行加密
type None struct {
}

func (crypto *None) Decode(data []byte) {

}

func (crypto *None) Encode(data []byte) {

}

func NewNoneCrypto() *None {
	return &None{}
}
