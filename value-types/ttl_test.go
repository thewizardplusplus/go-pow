package powValueTypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTTL(test *testing.T) {
	type args struct {
		rawValue time.Duration
	}

	for _, data := range []struct {
		name    string
		args    args
		want    TTL
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/positive",
			args: args{
				rawValue: 5*time.Minute + 23*time.Second,
			},
			want: TTL{
				rawValue: 5*time.Minute + 23*time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/zero",
			args: args{
				rawValue: 0,
			},
			want: TTL{
				rawValue: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: -(5*time.Minute + 23*time.Second),
			},
			want:    TTL{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewTTL(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestParseTTL(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name    string
		args    args
		want    TTL
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: "5m23s",
			},
			want: TTL{
				rawValue: 5*time.Minute + 23*time.Second,
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/invalid",
			args: args{
				rawValue: "invalid",
			},
			want:    TTL{},
			wantErr: assert.Error,
		},
		{
			name: "error/negative",
			args: args{
				rawValue: "-5m23s",
			},
			want:    TTL{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseTTL(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestTTL_ToDuration(test *testing.T) {
	type fields struct {
		rawValue time.Duration
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "success",
			fields: fields{
				rawValue: 5*time.Minute + 23*time.Second,
			},
			want: 5*time.Minute + 23*time.Second,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := TTL{
				rawValue: data.fields.rawValue,
			}
			got := value.ToDuration()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestTTL_ToString(test *testing.T) {
	type fields struct {
		rawValue time.Duration
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: 5*time.Minute + 23*time.Second,
			},
			want: "5m23s",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := TTL{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
