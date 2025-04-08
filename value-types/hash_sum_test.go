package powValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHashSum(test *testing.T) {
	type args struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name string
		args args
		want HashSum
	}{
		{
			name: "success/non-empty",
			args: args{
				rawValue: []byte("dummy"),
			},
			want: HashSum{
				rawValue: []byte("dummy"),
			},
		},
		{
			name: "success/empty",
			args: args{
				rawValue: []byte{},
			},
			want: HashSum{
				rawValue: []byte{},
			},
		},
		{
			name: "success/nil",
			args: args{
				rawValue: nil,
			},
			want: HashSum{
				rawValue: []byte{},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewHashSum(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashSum_Len(test *testing.T) {
	type fields struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				rawValue: []byte("dummy"),
			},
			want: 5,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HashSum{
				rawValue: data.fields.rawValue,
			}
			got := value.Len()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashSum_ToBytes(test *testing.T) {
	type fields struct {
		rawValue []byte
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "success",
			fields: fields{
				rawValue: []byte("dummy"),
			},
			want: []byte("dummy"),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HashSum{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBytes()

			assert.Equal(test, data.want, got)
		})
	}
}
