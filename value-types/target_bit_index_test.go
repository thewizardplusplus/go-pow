package powValueTypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTargetBitIndex(test *testing.T) {
	type args struct {
		rawValue int
	}

	for _, data := range []struct {
		name    string
		args    args
		want    TargetBitIndex
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/positive",
			args: args{
				rawValue: 23,
			},
			want: TargetBitIndex{
				rawValue: 23,
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/zero",
			args: args{
				rawValue: 0,
			},
			want: TargetBitIndex{
				rawValue: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: -23,
			},
			want:    TargetBitIndex{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewTargetBitIndex(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestTargetBitIndex_ToInt(test *testing.T) {
	type fields struct {
		rawValue int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "success",
			fields: fields{
				rawValue: 23,
			},
			want: 23,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := TargetBitIndex{
				rawValue: data.fields.rawValue,
			}
			got := value.ToInt()

			assert.Equal(test, data.want, got)
		})
	}
}
