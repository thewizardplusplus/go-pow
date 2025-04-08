package powValueTypes

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNonce(test *testing.T) {
	type args struct {
		rawValue *big.Int
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Nonce
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/positive",
			args: args{
				rawValue: big.NewInt(23),
			},
			want: Nonce{
				rawValue: big.NewInt(23),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/zero",
			args: args{
				rawValue: big.NewInt(0),
			},
			want: Nonce{
				rawValue: big.NewInt(0),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: big.NewInt(-23),
			},
			want:    Nonce{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewNonce(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestNewZeroNonce(test *testing.T) {
	for _, data := range []struct {
		name    string
		want    Nonce
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			want: Nonce{
				rawValue: big.NewInt(0),
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewZeroNonce()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestParseNonce(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Nonce
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: "23",
			},
			want: Nonce{
				rawValue: big.NewInt(23),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/invalid",
			args: args{
				rawValue: "invalid",
			},
			want:    Nonce{},
			wantErr: assert.Error,
		},
		{
			name: "error/negative",
			args: args{
				rawValue: "-23",
			},
			want:    Nonce{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseNonce(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestNonce_Incremented(test *testing.T) {
	type fields struct {
		rawValue *big.Int
	}

	for _, data := range []struct {
		name    string
		fields  fields
		want    Nonce
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				rawValue: big.NewInt(23),
			},
			want: Nonce{
				rawValue: big.NewInt(24),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				rawValue: big.NewInt(-23),
			},
			want:    Nonce{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Nonce{
				rawValue: data.fields.rawValue,
			}
			got, err := value.Incremented()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestNonce_ToBigInt(test *testing.T) {
	type fields struct {
		rawValue *big.Int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   *big.Int
	}{
		{
			name: "success",
			fields: fields{
				rawValue: big.NewInt(23),
			},
			want: big.NewInt(23),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Nonce{
				rawValue: data.fields.rawValue,
			}
			got := value.ToBigInt()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestNonce_ToString(test *testing.T) {
	type fields struct {
		rawValue *big.Int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: big.NewInt(23),
			},
			want: "23",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Nonce{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
