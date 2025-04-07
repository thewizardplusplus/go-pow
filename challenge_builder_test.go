package pow

import (
	"crypto/sha256"
	"net/url"
	"testing"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func TestChallengeBuilder_Build(test *testing.T) {
	for _, data := range []struct {
		name    string
		builder *ChallengeBuilder
		want    Challenge
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/all parameters",
			builder: NewChallengeBuilder().
				SetLeadingZeroCount(func() powValueTypes.LeadingZeroCount {
					value, err := powValueTypes.NewLeadingZeroCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetCreatedAt(func() powValueTypes.CreatedAt {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return value
				}()).
				SetResource(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})).
				SetPayload(powValueTypes.NewPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroCount.ToInt }}" +
						":{{ .Challenge.Payload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want: Challenge{
				leadingZeroCount: func() powValueTypes.LeadingZeroCount {
					value, err := powValueTypes.NewLeadingZeroCount(23)
					require.NoError(test, err)

					return value
				}(),
				createdAt: func() mo.Option[powValueTypes.CreatedAt] {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
				resource: mo.Some(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})),
				payload: powValueTypes.NewPayload("dummy"),
				hash:    powValueTypes.NewHash(sha256.New()),
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroCount.ToInt }}" +
						":{{ .Challenge.Payload.ToString }}" +
						":{{ .Nonce.ToString }}",
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only",
			builder: NewChallengeBuilder().
				SetLeadingZeroCount(func() powValueTypes.LeadingZeroCount {
					value, err := powValueTypes.NewLeadingZeroCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetPayload(powValueTypes.NewPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroCount.ToInt }}" +
						":{{ .Challenge.Payload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want: Challenge{
				leadingZeroCount: func() powValueTypes.LeadingZeroCount {
					value, err := powValueTypes.NewLeadingZeroCount(23)
					require.NoError(test, err)

					return value
				}(),
				createdAt: mo.None[powValueTypes.CreatedAt](),
				resource:  mo.None[powValueTypes.Resource](),
				payload:   powValueTypes.NewPayload("dummy"),
				hash:      powValueTypes.NewHash(sha256.New()),
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroCount.ToInt }}" +
						":{{ .Challenge.Payload.ToString }}" +
						":{{ .Nonce.ToString }}",
				),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error/without parameters",
			builder: NewChallengeBuilder(),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/leading zero count is too large",
			builder: NewChallengeBuilder().
				SetLeadingZeroCount(func() powValueTypes.LeadingZeroCount {
					value, err := powValueTypes.NewLeadingZeroCount(1000)
					require.NoError(test, err)

					return value
				}()).
				SetPayload(powValueTypes.NewPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroCount.ToInt }}" +
						":{{ .Challenge.Payload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got, err := data.builder.Build()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
