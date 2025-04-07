package powValueTypes

import (
	"hash"
	"reflect"
)

const (
	bitsPerByte = 8
)

type Hash struct {
	rawValue hash.Hash
}

func NewHash(rawValue hash.Hash) Hash {
	return Hash{
		rawValue: rawValue,
	}
}

func (value Hash) Name() string {
	// don't use `reflect.Type.Name()`, as implementations of `hash.Hash`
	// are most often private
	return reflect.TypeOf(value.rawValue).String()
}

func (value Hash) SizeInBytes() int {
	return value.rawValue.Size()
}

func (value Hash) SizeInBits() int {
	return value.rawValue.Size() * bitsPerByte
}

func (value Hash) ApplyTo(data string) []byte {
	value.rawValue.Reset()
	value.rawValue.Write([]byte(data))
	return value.rawValue.Sum(nil)
}

func (value Hash) ToHash() hash.Hash {
	return value.rawValue
}
