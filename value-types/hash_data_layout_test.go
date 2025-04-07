package powValueTypes

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestNewHashDataLayout(test *testing.T) {
	type args struct {
		rawValue *template.Template
	}

	for _, data := range []struct {
		name string
		args args
		want HashDataLayout
	}{
		{
			name: "success",
			args: args{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			want: HashDataLayout{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewHashDataLayout(data.args.rawValue)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestParseHashDataLayout(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name    string
		args    args
		want    HashDataLayout
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: "dummy {{ .Dummy }}",
			},
			want: HashDataLayout{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				rawValue: "dummy {{ .Dummy",
			},
			want:    HashDataLayout{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := ParseHashDataLayout(data.args.rawValue)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestMustParseHashDataLayout(test *testing.T) {
	type args struct {
		rawValue string
	}

	for _, data := range []struct {
		name      string
		args      args
		want      HashDataLayout
		wantPanic assert.PanicAssertionFunc
	}{
		{
			name: "success",
			args: args{
				rawValue: "dummy {{ .Dummy }}",
			},
			want: HashDataLayout{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			wantPanic: assert.NotPanics,
		},
		{
			name: "error",
			args: args{
				rawValue: "dummy {{ .Dummy",
			},
			want:      HashDataLayout{},
			wantPanic: assert.Panics,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var got HashDataLayout
			data.wantPanic(test, func() {
				got = MustParseHashDataLayout(data.args.rawValue)
			})

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashDataLayout_Execute(test *testing.T) {
	type fields struct {
		rawValue *template.Template
	}
	type args struct {
		data any
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/numeric data",
			fields: fields{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			args: args{
				data: struct {
					Dummy int
				}{
					Dummy: 23,
				},
			},
			want:    "dummy 23",
			wantErr: assert.NoError,
		},
		{
			name: "success/string data",
			fields: fields{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			args: args{
				data: struct {
					Dummy string
				}{
					Dummy: "suffix",
				},
			},
			want:    "dummy suffix",
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				rawValue: template.Must(template.New("").Parse("dummy {{ .DummyOne }}")),
			},
			args: args{
				data: struct {
					DummyTwo int
				}{
					DummyTwo: 23,
				},
			},
			want:    "",
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HashDataLayout{
				rawValue: data.fields.rawValue,
			}
			got, err := value.Execute(data.args.data)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestHashDataLayout_ToTemplate(test *testing.T) {
	type fields struct {
		rawValue *template.Template
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   *template.Template
	}{
		{
			name: "success",
			fields: fields{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			want: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HashDataLayout{
				rawValue: data.fields.rawValue,
			}
			got := value.ToTemplate()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashDataLayout_ToString(test *testing.T) {
	type fields struct {
		rawValue *template.Template
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success",
			fields: fields{
				rawValue: template.Must(template.New("").Parse("dummy {{ .Dummy }}")),
			},
			want: "dummy {{.Dummy}}",
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := HashDataLayout{
				rawValue: data.fields.rawValue,
			}
			got := value.ToString()

			assert.Equal(test, data.want, got)
		})
	}
}
