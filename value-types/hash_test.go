package powValueTypes

import (
	"crypto/sha256"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHash(test *testing.T) {
	type args struct {
		rawValue hash.Hash
	}

	for _, data := range []struct {
		name string
		args args
		want Hash
	}{
		{
			name: "success",
			args: args{
				rawValue: sha256.New(),
			},
			want: Hash{
				rawValue: sha256.New(),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewHash(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHash_Name(test *testing.T) {
	type fields struct {
		rawValue hash.Hash
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: sha256.New(),
			},
			want: "*sha256.digest",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Hash{
				rawValue: data.fields.rawValue,
			}
			got := value.Name()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHash_SizeInBytes(test *testing.T) {
	type fields struct {
		rawValue hash.Hash
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				rawValue: sha256.New(),
			},
			want: 32,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Hash{
				rawValue: data.fields.rawValue,
			}
			got := value.SizeInBytes()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHash_SizeInBits(test *testing.T) {
	type fields struct {
		rawValue hash.Hash
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				rawValue: sha256.New(),
			},
			want: 256,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Hash{
				rawValue: data.fields.rawValue,
			}
			got := value.SizeInBits()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHash_ApplyTo(test *testing.T) {
	type fields struct {
		rawValue hash.Hash
	}
	type args struct {
		data string
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "success/from scratch",
			fields: fields{
				rawValue: sha256.New(),
			},
			args: args{
				data: "dummy",
			},
			want: []byte{
				0xb5, 0xa2, 0xc9, 0x62, 0x50, 0x61, 0x23, 0x66,
				0xea, 0x27, 0x2f, 0xfa, 0xc6, 0xd9, 0x74, 0x4a,
				0xaf, 0x4b, 0x45, 0xaa, 0xcd, 0x96, 0xaa, 0x7c,
				0xfc, 0xb9, 0x31, 0xee, 0x3b, 0x55, 0x82, 0x59,
			},
		},
		{
			name: "success/not from scratch",
			fields: fields{
				rawValue: func() hash.Hash {
					hash := sha256.New()
					hash.Write([]byte("prefix"))

					return hash
				}(),
			},
			args: args{
				data: "dummy",
			},
			want: []byte{
				0xb5, 0xa2, 0xc9, 0x62, 0x50, 0x61, 0x23, 0x66,
				0xea, 0x27, 0x2f, 0xfa, 0xc6, 0xd9, 0x74, 0x4a,
				0xaf, 0x4b, 0x45, 0xaa, 0xcd, 0x96, 0xaa, 0x7c,
				0xfc, 0xb9, 0x31, 0xee, 0x3b, 0x55, 0x82, 0x59,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Hash{
				rawValue: data.fields.rawValue,
			}
			got := value.ApplyTo(data.args.data)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHash_ToHash(test *testing.T) {
	type fields struct {
		rawValue hash.Hash
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   hash.Hash
	}{
		{
			name: "success",
			fields: fields{
				rawValue: sha256.New(),
			},
			want: sha256.New(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Hash{
				rawValue: data.fields.rawValue,
			}
			got := value.ToHash()

			assert.Equal(test, data.want, got)
		})
	}
}
