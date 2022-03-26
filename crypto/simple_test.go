package crypto

import "testing"

func TestSimple(t *testing.T) {
	crypto := NewSimpleCrypto("2wYzDgXzAM/XvhRUNIr4P5LlRvXNbn0CrozAH2N0yxX9aXdbLnaEHe21T6jcbC25A6vhk714cCizJwzqgHlD6+me03OIoF/HQJh7Omj0XaEIWBMNcuhr0GqZSrbyo5G4C9EBt/GGEbFnyORJXrtXClOW1dkvbx7vj5RhGBqC3mAhzoklf/6/8EErONRmZd3/RdiVwcmsHOMJgcVtpVI+xKmL+6KDJJwmvOykTXxZMZAqTLQgFjlxPVywIxLfMIVa2sOXIuedSLLuLJvWdfYy5jw1yuI2wimqYpqORDc796e6rRevh/l6+gfgpksEGVZVjRAPfsYbQsxRn2RH0k78UA==")
	content := []byte("this is a cat")
	t.Logf("original:%v\n", content)
	crypto.Encode(content)
	t.Logf("encode:%v\n", content)
	crypto.Decode(content)
	t.Logf("decode:%v\n", content)
}
