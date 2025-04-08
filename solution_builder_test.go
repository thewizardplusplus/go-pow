package pow

import (
	"crypto/sha256"
	"math/big"
	"testing"

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
			name: "success",
			builder: NewSolutionBuilder().
				SetChallenge(Challenge{
					leadingZeroCount: func() powValueTypes.LeadingZeroCount {
						value, err := powValueTypes.NewLeadingZeroCount(23)
						require.NoError(test, err)

						return value
					}(),
					payload: powValueTypes.NewPayload("dummy"),
					hash:    powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroCount.ToInt }}" +
							":{{ .Challenge.Payload.ToString }}" +
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
					leadingZeroCount: func() powValueTypes.LeadingZeroCount {
						value, err := powValueTypes.NewLeadingZeroCount(23)
						require.NoError(test, err)

						return value
					}(),
					payload: powValueTypes.NewPayload("dummy"),
					hash:    powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroCount.ToInt }}" +
							":{{ .Challenge.Payload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				},
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			builder: NewSolutionBuilder(),
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
