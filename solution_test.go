package pow

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/samber/mo"
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
			},
			want: Challenge{
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
		hashSum mo.Option[powValueTypes.HashSum]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[powValueTypes.HashSum]
	}{
		{
			name: "success/is present",
			fields: fields{
				hashSum: mo.Some(powValueTypes.NewHashSum([]byte("dummy"))),
			},
			want: mo.Some(powValueTypes.NewHashSum([]byte("dummy"))),
		},
		{
			name: "success/is absent",
			fields: fields{
				hashSum: mo.None[powValueTypes.HashSum](),
			},
			want: mo.None[powValueTypes.HashSum](),
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

func TestSolution_Verify(test *testing.T) {
	type fields struct {
		challenge Challenge
		nonce     powValueTypes.Nonce
		hashSum   mo.Option[powValueTypes.HashSum]
	}

	for _, data := range []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/hash sum is present",
			fields: fields{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
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
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.Some(powValueTypes.NewHashSum([]byte{
					0x00, 0x5d, 0x37, 0x2c, 0x56, 0xe6, 0xc6, 0xb5,
					0x2a, 0xd4, 0xa8, 0x32, 0x56, 0x54, 0x69, 0x2e,
					0xc9, 0xaa, 0x3a, 0xf5, 0xf7, 0x30, 0x21, 0x74,
					0x8b, 0xc3, 0xfd, 0xb1, 0x24, 0xae, 0x9b, 0x20,
				})),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success/hash sum is absent",
			fields: fields{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
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
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.None[powValueTypes.HashSum](),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to get the target bit index",
			fields: fields{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(1000)
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
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.None[powValueTypes.HashSum](),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to execute the hash data layout",
			fields: fields{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash:              powValueTypes.NewHash(sha256.New()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"dummy {{ .Dummy }}",
					),
				},
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.None[powValueTypes.HashSum](),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/hash sums don't match",
			fields: fields{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
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
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.Some(powValueTypes.NewHashSum([]byte("dummy"))),
			},
			wantErr: assert.Error,
		},
		{
			name: "error/hash sum doesn't fit the target",
			fields: fields{
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
					value, err := powValueTypes.NewNonce(big.NewInt(37))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.Some(powValueTypes.NewHashSum([]byte{
					0x2d, 0x55, 0x8a, 0x78, 0xdf, 0x38, 0xa3, 0xe4,
					0x1c, 0x3f, 0x53, 0x24, 0xeb, 0x32, 0xaa, 0x31,
					0x3b, 0x3f, 0xa7, 0xc3, 0xb4, 0xd3, 0xe8, 0x2f,
					0x2b, 0x5d, 0x98, 0x96, 0xd1, 0xa2, 0x36, 0x34,
				})),
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Solution{
				challenge: data.fields.challenge,
				nonce:     data.fields.nonce,
				hashSum:   data.fields.hashSum,
			}
			err := entity.Verify()

			data.wantErr(test, err)
		})
	}
}
