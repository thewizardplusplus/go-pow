package powValueTypes

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(test *testing.T) {
	type args struct {
		rawValue *url.URL
	}

	for _, data := range []struct {
		name string
		args args
		want Resource
	}{
		{
			name: "success",
			args: args{
				rawValue: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				},
			},
			want: Resource{
				rawValue: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewResource(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestParseResource(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name    string
		args    args
		want    Resource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/https",
			args: args{
				rawValue: "https://example.com/",
			},
			want: Resource{
				rawValue: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/mailto",
			args: args{
				rawValue: "mailto:someone@example.com",
			},
			want: Resource{
				rawValue: &url.URL{
					Scheme: "mailto",
					Opaque: "someone@example.com",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: ":",
			},
			want:    Resource{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseResource(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestResource_ToURL(test *testing.T) {
	type fields struct {
		rawValue *url.URL
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   *url.URL
	}{
		{
			name: "success",
			fields: fields{
				rawValue: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				},
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/",
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Resource{
				rawValue: data.fields.rawValue,
			}
			got := value.ToURL()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestResource_ToString(test *testing.T) {
	type fields struct {
		rawValue *url.URL
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: &url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				},
			},
			want: "https://example.com/",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Resource{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
