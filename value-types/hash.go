package powValueTypes

import (
	"errors"
	"hash"
	"reflect"

	"github.com/samber/mo"
)

const (
	bitsPerByte = 8
)

type Hash struct {
	rawValue hash.Hash
	name     mo.Option[string]
}

func NewHash(rawValue hash.Hash) Hash {
	return Hash{
		rawValue: rawValue,
	}
}

func NewHashWithName(rawValue hash.Hash, name string) (Hash, error) {
	if name == "" {
		return Hash{}, errors.New("hash name cannot be empty")
	}

	value := Hash{
		rawValue: rawValue,
		name:     mo.Some(name),
	}
	return value, nil
}

func (value Hash) Name() string {
	name, isPresent := value.name.Get()
	if isPresent {
		return name
	}

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

func (value Hash) ApplyTo(data string) HashSum {
	value.rawValue.Reset()
	value.rawValue.Write([]byte(data))
	return NewHashSum(value.rawValue.Sum(nil))
}

func (value Hash) ToHash() hash.Hash {
	return value.rawValue
}
