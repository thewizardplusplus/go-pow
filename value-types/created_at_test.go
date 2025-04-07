package powValueTypes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCreatedAt(test *testing.T) {
	type args struct {
		rawValue time.Time
	}

	for _, data := range []struct {
		name    string
		args    args
		want    CreatedAt
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
			},
			want: CreatedAt{
				rawValue: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: time.Time{},
			},
			want:    CreatedAt{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := NewCreatedAt(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestParseCreatedAt(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name    string
		args    args
		want    CreatedAt
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: "2000-01-02T03:04:05.000000006Z",
			},
			want: CreatedAt{
				rawValue: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/invalid",
			args: args{
				rawValue: "invalid",
			},
			want:    CreatedAt{},
			wantErr: assert.Error,
		},
		{
			name: "error/zero",
			args: args{
				rawValue: "0001-01-01T00:00:00Z",
			},
			want:    CreatedAt{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseCreatedAt(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestCreatedAt_ToTime(test *testing.T) {
	type fields struct {
		rawValue time.Time
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "success",
			fields: fields{
				rawValue: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
			},
			want: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := CreatedAt{
				rawValue: data.fields.rawValue,
			}
			got := value.ToTime()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestCreatedAt_ToString(test *testing.T) {
	type fields struct {
		rawValue time.Time
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
			},
			want: "2000-01-02T03:04:05.000000006Z",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := CreatedAt{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
