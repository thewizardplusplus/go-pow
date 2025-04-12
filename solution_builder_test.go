package pow

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func TestSolutionBuilder_Build(test *testing.T) {
	for _, data := range []struct {
		name    string
		builder *SolutionBuilder
		want    Solution
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/all parameters",
			builder: NewSolutionBuilder().
				SetChallenge(Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(23)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				}).
				SetNonce(func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}()).
				SetHashSum(powValueTypes.NewHashSum(bytes.Repeat([]byte("0"), 32))),
			want: Solution{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(23)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				},
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.Some(powValueTypes.NewHashSum(bytes.Repeat([]byte("0"), 32))),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/required parameters only",
			builder: NewSolutionBuilder().
				SetChallenge(Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(23)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				}).
				SetNonce(func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}()),
			want: Solution{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(23)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				},
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.None[powValueTypes.HashSum](),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error/without parameters",
			builder: NewSolutionBuilder(),
			want:    Solution{},
			wantErr: assert.Error,
		},
		{
			name: "error/invalid hash sum length",
			builder: NewSolutionBuilder().
				SetChallenge(Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(23)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				}).
				SetNonce(func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}()).
				SetHashSum(powValueTypes.NewHashSum([]byte("dummy"))),
			want:    Solution{},
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
