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
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
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
				SetTTL(func() powValueTypes.TTL {
					value, err := powValueTypes.NewTTL(5*time.Minute + 23*time.Second)
					require.NoError(test, err)

					return value
				}()).
				SetResource(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want: Challenge{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
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
				ttl: func() mo.Option[powValueTypes.TTL] {
					value, err := powValueTypes.NewTTL(5*time.Minute + 23*time.Second)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
				resource: mo.Some(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})),
				serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
				hash:              powValueTypes.NewHash(sha256.New()),
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only/leading zero bit count is specified",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want: Challenge{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}(),
				createdAt:         mo.None[powValueTypes.CreatedAt](),
				resource:          mo.None[powValueTypes.Resource](),
				serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
				hash:              powValueTypes.NewHash(sha256.New()),
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only/target bit index is specified",
			builder: NewChallengeBuilder().
				SetTargetBitIndex(func() powValueTypes.TargetBitIndex {
					value, err := powValueTypes.NewTargetBitIndex(233)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want: Challenge{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}(),
				createdAt:         mo.None[powValueTypes.CreatedAt](),
				resource:          mo.None[powValueTypes.Resource](),
				serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
				hash:              powValueTypes.NewHash(sha256.New()),
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
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
			name: "error/" +
				"leading zero bit count and target bit index " +
				"are specified at the same time",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetTargetBitIndex(func() powValueTypes.TargetBitIndex {
					value, err := powValueTypes.NewTargetBitIndex(233)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/" +
				"`CreatedAt` timestamp and TTL " +
				"should either both be specified or both omitted/" +
				"`CreatedAt` timestamp is specified",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
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
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/" +
				"`CreatedAt` timestamp and TTL " +
				"should either both be specified or both omitted/" +
				"TTL is specified",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetTTL(func() powValueTypes.TTL {
					value, err := powValueTypes.NewTTL(5*time.Minute + 23*time.Second)
					require.NoError(test, err)

					return value
				}()).
				SetResource(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/leading zero bit count is too large",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(1000)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/target bit index is too large",
			builder: NewChallengeBuilder().
				SetTargetBitIndex(func() powValueTypes.TargetBitIndex {
					value, err := powValueTypes.NewTargetBitIndex(1000)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				)),
			want:    Challenge{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to check the hash data layout",
			builder: NewChallengeBuilder().
				SetLeadingZeroBitCount(func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}()).
				SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
				SetHash(powValueTypes.NewHash(sha256.New())).
				SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
					"dummy {{ .Dummy }}",
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
