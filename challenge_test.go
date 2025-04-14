package pow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"hash"
	"math/big"
	"net/url"
	"testing"
	"testing/iotest"
	"time"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func TestChallenge_LeadingZeroBitCount(test *testing.T) {
	type fields struct {
		leadingZeroBitCount powValueTypes.LeadingZeroBitCount
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.LeadingZeroBitCount
	}{
		{
			name: "success",
			fields: fields{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}(),
			},
			want: func() powValueTypes.LeadingZeroBitCount {
				value, err := powValueTypes.NewLeadingZeroBitCount(23)
				require.NoError(test, err)

				return value
			}(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				leadingZeroBitCount: data.fields.leadingZeroBitCount,
			}
			got := entity.LeadingZeroBitCount()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_TargetBitIndex(test *testing.T) {
	type fields struct {
		leadingZeroBitCount powValueTypes.LeadingZeroBitCount
		hash                powValueTypes.Hash
	}

	for _, data := range []struct {
		name    string
		fields  fields
		want    powValueTypes.TargetBitIndex
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/regular leading zero bit count",
			fields: fields{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(23)
					require.NoError(test, err)

					return value
				}(),
				hash: powValueTypes.NewHash(sha256.New()),
			},
			want: func() powValueTypes.TargetBitIndex {
				value, err := powValueTypes.NewTargetBitIndex(233)
				require.NoError(test, err)

				return value
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success/minimal leading zero bit count",
			fields: fields{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(0)
					require.NoError(test, err)

					return value
				}(),
				hash: powValueTypes.NewHash(sha256.New()),
			},
			want: func() powValueTypes.TargetBitIndex {
				value, err := powValueTypes.NewTargetBitIndex(256)
				require.NoError(test, err)

				return value
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success/maximal leading zero bit count",
			fields: fields{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(256)
					require.NoError(test, err)

					return value
				}(),
				hash: powValueTypes.NewHash(sha256.New()),
			},
			want: func() powValueTypes.TargetBitIndex {
				value, err := powValueTypes.NewTargetBitIndex(0)
				require.NoError(test, err)

				return value
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
					value, err := powValueTypes.NewLeadingZeroBitCount(1000)
					require.NoError(test, err)

					return value
				}(),
				hash: powValueTypes.NewHash(sha256.New()),
			},
			want:    powValueTypes.TargetBitIndex{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				leadingZeroBitCount: data.fields.leadingZeroBitCount,
				hash:                data.fields.hash,
			}
			got, err := entity.TargetBitIndex()

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}

func TestChallenge_CreatedAt(test *testing.T) {
	type fields struct {
		createdAt mo.Option[powValueTypes.CreatedAt]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[powValueTypes.CreatedAt]
	}{
		{
			name: "success/is present",
			fields: fields{
				createdAt: func() mo.Option[powValueTypes.CreatedAt] {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
			},
			want: func() mo.Option[powValueTypes.CreatedAt] {
				value, err := powValueTypes.NewCreatedAt(
					time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
				)
				require.NoError(test, err)

				return mo.Some(value)
			}(),
		},
		{
			name: "success/is absent",
			fields: fields{
				createdAt: mo.None[powValueTypes.CreatedAt](),
			},
			want: mo.None[powValueTypes.CreatedAt](),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				createdAt: data.fields.createdAt,
			}
			got := entity.CreatedAt()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_TTL(test *testing.T) {
	type fields struct {
		ttl mo.Option[powValueTypes.TTL]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[powValueTypes.TTL]
	}{
		{
			name: "success/is present",
			fields: fields{
				ttl: func() mo.Option[powValueTypes.TTL] {
					value, err := powValueTypes.NewTTL(5*time.Minute + 23*time.Second)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
			},
			want: func() mo.Option[powValueTypes.TTL] {
				value, err := powValueTypes.NewTTL(5*time.Minute + 23*time.Second)
				require.NoError(test, err)

				return mo.Some(value)
			}(),
		},
		{
			name: "success/is absent",
			fields: fields{
				ttl: mo.None[powValueTypes.TTL](),
			},
			want: mo.None[powValueTypes.TTL](),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				ttl: data.fields.ttl,
			}
			got := entity.TTL()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_IsAlive(test *testing.T) {
	type fields struct {
		createdAt mo.Option[powValueTypes.CreatedAt]
		ttl       mo.Option[powValueTypes.TTL]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		{
			name: "success/is alive/within the TTL",
			fields: fields{
				createdAt: func() mo.Option[powValueTypes.CreatedAt] {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
				ttl: func() mo.Option[powValueTypes.TTL] {
					value, err := powValueTypes.NewTTL(100 * 365 * 24 * time.Hour)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
			},
			want: assert.True,
		},
		{
			name: "success/is alive/`CreatedAt` timestamp isn't specified",
			fields: fields{
				createdAt: mo.None[powValueTypes.CreatedAt](),
				ttl: func() mo.Option[powValueTypes.TTL] {
					value, err := powValueTypes.NewTTL(100 * 365 * 24 * time.Hour)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
			},
			want: assert.True,
		},
		{
			name: "success/is alive/TTL isn't specified",
			fields: fields{
				createdAt: func() mo.Option[powValueTypes.CreatedAt] {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
				ttl: mo.None[powValueTypes.TTL](),
			},
			want: assert.True,
		},
		{
			name: "success/is dead",
			fields: fields{
				createdAt: func() mo.Option[powValueTypes.CreatedAt] {
					value, err := powValueTypes.NewCreatedAt(
						time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC),
					)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
				ttl: func() mo.Option[powValueTypes.TTL] {
					value, err := powValueTypes.NewTTL(time.Second)
					require.NoError(test, err)

					return mo.Some(value)
				}(),
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				createdAt: data.fields.createdAt,
				ttl:       data.fields.ttl,
			}
			got := entity.IsAlive()

			data.want(test, got)
		})
	}
}

func TestChallenge_Resource(test *testing.T) {
	type fields struct {
		resource mo.Option[powValueTypes.Resource]
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   mo.Option[powValueTypes.Resource]
	}{
		{
			name: "success/is present",
			fields: fields{
				resource: mo.Some(powValueTypes.NewResource(&url.URL{
					Scheme: "https",
					Host:   "example.com",
					Path:   "/",
				})),
			},
			want: mo.Some(powValueTypes.NewResource(&url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/",
			})),
		},
		{
			name: "success/is absent",
			fields: fields{
				resource: mo.None[powValueTypes.Resource](),
			},
			want: mo.None[powValueTypes.Resource](),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				resource: data.fields.resource,
			}
			got := entity.Resource()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_SerializedPayload(test *testing.T) {
	type fields struct {
		serializedPayload powValueTypes.SerializedPayload
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.SerializedPayload
	}{
		{
			name: "success",
			fields: fields{
				serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
			},
			want: powValueTypes.NewSerializedPayload("dummy"),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				serializedPayload: data.fields.serializedPayload,
			}
			got := entity.SerializedPayload()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_Hash(test *testing.T) {
	type fields struct {
		hash powValueTypes.Hash
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.Hash
	}{
		{
			name: "success",
			fields: fields{
				hash: powValueTypes.NewHash(sha256.New()),
			},
			want: powValueTypes.NewHash(sha256.New()),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				hash: data.fields.hash,
			}
			got := entity.Hash()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_HashDataLayout(test *testing.T) {
	type fields struct {
		hashDataLayout powValueTypes.HashDataLayout
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   powValueTypes.HashDataLayout
	}{
		{
			name: "success",
			fields: fields{
				hashDataLayout: powValueTypes.MustParseHashDataLayout(
					"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
						":{{ .Challenge.SerializedPayload.ToString }}" +
						":{{ .Nonce.ToString }}",
				),
			},
			want: powValueTypes.MustParseHashDataLayout(
				"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
					":{{ .Challenge.SerializedPayload.ToString }}" +
					":{{ .Nonce.ToString }}",
			),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				hashDataLayout: data.fields.hashDataLayout,
			}
			got := entity.HashDataLayout()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestChallenge_Solve(test *testing.T) {
	type fields struct {
		leadingZeroBitCount powValueTypes.LeadingZeroBitCount
		serializedPayload   powValueTypes.SerializedPayload
		hash                powValueTypes.Hash
		hashDataLayout      powValueTypes.HashDataLayout
	}
	type args struct {
		ctx    context.Context
		params SolveParams
	}

	for _, data := range []struct {
		name    string
		fields  fields
		args    args
		want    Solution
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success/zero initial nonce",
			fields: fields{
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
			args: args{
				ctx:    context.Background(),
				params: SolveParams{},
			},
			want: Solution{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash: powValueTypes.NewHash(func() hash.Hash {
						hash := sha256.New()
						hash.Write([]byte("5:dummy:37"))

						return hash
					}()),
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
			name: "success/random initial nonce",
			fields: fields{
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
			args: args{
				ctx: context.Background(),
				params: SolveParams{
					RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
						RandomReader: bytes.NewReader([]byte("dummy")),
						MinRawValue:  big.NewInt(123),
						MaxRawValue:  big.NewInt(142),
					}),
				},
			},
			want: Solution{
				challenge: Challenge{
					leadingZeroBitCount: func() powValueTypes.LeadingZeroBitCount {
						value, err := powValueTypes.NewLeadingZeroBitCount(5)
						require.NoError(test, err)

						return value
					}(),
					serializedPayload: powValueTypes.NewSerializedPayload("dummy"),
					hash: powValueTypes.NewHash(func() hash.Hash {
						hash := sha256.New()
						hash.Write([]byte("5:dummy:129"))

						return hash
					}()),
					hashDataLayout: powValueTypes.MustParseHashDataLayout(
						"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
							":{{ .Challenge.SerializedPayload.ToString }}" +
							":{{ .Nonce.ToString }}",
					),
				},
				nonce: func() powValueTypes.Nonce {
					value, err := powValueTypes.NewNonce(big.NewInt(129))
					require.NoError(test, err)

					return value
				}(),
				hashSum: mo.Some(powValueTypes.NewHashSum([]byte{
					0x01, 0x2b, 0x29, 0x61, 0x88, 0x31, 0x8c, 0xa1,
					0xe0, 0x93, 0x88, 0x91, 0x93, 0x32, 0x14, 0x48,
					0xd5, 0xd7, 0x11, 0x49, 0x46, 0xe6, 0x68, 0x1c,
					0x07, 0x75, 0xdb, 0x4c, 0xdc, 0x4e, 0x76, 0xb3,
				})),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error/unable to generate the random initial nonce",
			fields: fields{
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
			args: args{
				ctx: context.Background(),
				params: SolveParams{
					RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
						RandomReader: iotest.ErrReader(iotest.ErrTimeout),
						MinRawValue:  big.NewInt(123),
						MaxRawValue:  big.NewInt(142),
					}),
				},
			},
			want:    Solution{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to get the target bit index",
			fields: fields{
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
			args: args{
				ctx:    context.Background(),
				params: SolveParams{},
			},
			want:    Solution{},
			wantErr: assert.Error,
		},
		{
			name: "error/context is done",
			fields: fields{
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
			args: args{
				ctx: func() context.Context {
					ctx, ctxCancel := context.WithCancel(context.Background())
					ctxCancel()

					return ctx
				}(),
				params: SolveParams{},
			},
			want:    Solution{},
			wantErr: assert.Error,
		},
		{
			name: "error/maximal attempt count is exceeded",
			fields: fields{
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
			args: args{
				ctx: context.Background(),
				params: SolveParams{
					MaxAttemptCount: mo.Some(23),
				},
			},
			want:    Solution{},
			wantErr: assert.Error,
		},
		{
			name: "error/unable to execute the hash data layout",
			fields: fields{
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
			args: args{
				ctx:    context.Background(),
				params: SolveParams{},
			},
			want:    Solution{},
			wantErr: assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			entity := Challenge{
				leadingZeroBitCount: data.fields.leadingZeroBitCount,
				serializedPayload:   data.fields.serializedPayload,
				hash:                data.fields.hash,
				hashDataLayout:      data.fields.hashDataLayout,
			}
			got, err := entity.Solve(data.args.ctx, data.args.params)

			assert.Equal(test, data.want, got)
			data.wantErr(test, err)
		})
	}
}
