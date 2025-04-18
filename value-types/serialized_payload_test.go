package powValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSerializedPayload(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name string
		args args
		want SerializedPayload
	}{
		{
			name: "success/non-empty",
			args: args{
				rawValue: "dummy",
			},
			want: SerializedPayload{
				rawValue: "dummy",
			},
		},
		{
			name: "success/empty",
			args: args{
				rawValue: "",
			},
			want: SerializedPayload{
				rawValue: "",
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewSerializedPayload(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestSerializedPayload_ToString(test *testing.T) {
	type fields struct {
		rawValue string
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: "dummy",
			},
			want: "dummy",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := SerializedPayload{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
