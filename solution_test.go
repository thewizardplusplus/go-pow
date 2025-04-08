package pow

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func TestSolution_Challenge(test *testing.T) {
	type fields struct {
		challenge Challenge
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   Challenge
	}{
		{
			name: "success",
			fields: fields{
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
			},
			want: Challenge{
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
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Solution{
				challenge: data.fields.challenge,
			}
			got := entity.Challenge()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestSolution_Nonce(test *testing.T) {
	type fields struct {
		nonce powValueTypes.Nonce
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.Nonce
	}{
		{
			name: "success",
			fields: fields{
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(23))
					require.NoError(test, err)

					return value
				}(),
			},
			want: func() powValueTypes.Nonce {
				value, err := powValueTypes.NewNonce(big.NewInt(23))
				require.NoError(test, err)

				return value
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Solution{
				nonce: data.fields.nonce,
			}
			got := entity.Nonce()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestSolution_HashSum(test *testing.T) {
	type fields struct {
		hashSum powValueTypes.HashSum
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.HashSum
	}{
		{
			name: "success",
			fields: fields{
				hashSum: powValueTypes.NewHashSum([]byte("dummy")),
			},
			want: powValueTypes.NewHashSum([]byte("dummy")),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Solution{
				hashSum: data.fields.hashSum,
			}
			got := entity.HashSum()

			assert.Equal(test, data.want, got)
		})
	}
}
