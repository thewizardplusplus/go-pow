package pow_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/samber/mo"
	pow "github.com/thewizardplusplus/go-pow"
	powErrors "github.com/thewizardplusplus/go-pow/errors"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func Example_minimal() {
	leadingZeroBitCount, err := powValueTypes.NewLeadingZeroBitCount(5)
	if err != nil {
		log.Fatalf("unable to construct the leading zero bit count: %s", err)
	}

	namedHash, err := powValueTypes.NewHashWithName(sha256.New(), "SHA-256")
	if err != nil {
		log.Fatalf("unable to construct the hash: %s", err)
	}

	challenge, err := pow.NewChallengeBuilder().
		SetLeadingZeroBitCount(leadingZeroBitCount).
		SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
		SetHash(namedHash).
		SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
			"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
				":{{ .Challenge.SerializedPayload.ToString }}" +
				":{{ .Nonce.ToString }}",
		)).
		Build()
	if err != nil {
		log.Fatalf("unable to build the challenge: %s", err)
	}

	fmt.Printf("challenge: %v\n", []string{
		strconv.Itoa(challenge.LeadingZeroBitCount().ToInt()),
		challenge.SerializedPayload().ToString(),
		challenge.Hash().Name(),
		challenge.HashDataLayout().ToString(),
	})

	solution, err := challenge.Solve(context.Background(), pow.SolveParams{})
	if err != nil {
		log.Fatalf("unable to solve the challenge: %s", err)
	}

	fmt.Printf("nonce: %s\n", solution.Nonce().ToString())
	fmt.Printf("hash sum: %x\n", solution.HashSum().OrEmpty().ToBytes())

	if err := solution.Verify(); err != nil {
		log.Fatalf("unable to verify the solution: %s", err)
	}

	fmt.Print("verification: OK\n")

	// Output:
	// challenge: [5 dummy SHA-256 {{.Challenge.LeadingZeroBitCount.ToInt}}:{{.Challenge.SerializedPayload.ToString}}:{{.Nonce.ToString}}]
	// nonce: 37
	// hash sum: 005d372c56e6c6b52ad4a8325654692ec9aa3af5f73021748bc3fdb124ae9b20
	// verification: OK
}

func Example_full() {
	leadingZeroBitCount, err := powValueTypes.NewLeadingZeroBitCount(5)
	if err != nil {
		log.Fatalf("unable to construct the leading zero bit count: %s", err)
	}

	rawCreatedAt := time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC)
	createdAt, err := powValueTypes.NewCreatedAt(rawCreatedAt)
	if err != nil {
		log.Fatalf("unable to parse the `CreatedAt` timestamp: %s", err)
	}

	ttl, err := powValueTypes.NewTTL(100 * 365 * 24 * time.Hour)
	if err != nil {
		log.Fatalf("unable to parse the TTL: %s", err)
	}

	namedHash, err := powValueTypes.NewHashWithName(sha256.New(), "SHA-256")
	if err != nil {
		log.Fatalf("unable to construct the hash: %s", err)
	}

	challenge, err := pow.NewChallengeBuilder().
		SetLeadingZeroBitCount(leadingZeroBitCount).
		SetCreatedAt(createdAt).
		SetTTL(ttl).
		SetResource(powValueTypes.NewResource(&url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/",
		})).
		SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
		SetHash(namedHash).
		SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
			"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
				":{{ .Challenge.CreatedAt.MustGet.ToString }}" +
				":{{ .Challenge.TTL.MustGet.ToString }}" +
				":{{ .Challenge.Resource.MustGet.ToString }}" +
				":{{ .Challenge.SerializedPayload.ToString }}" +
				":{{ .Challenge.Hash.Name }}" +
				":{{ .Challenge.HashDataLayout.ToString }}" +
				":{{ .Nonce.ToString }}",
		)).
		Build()
	if err != nil {
		log.Fatalf("unable to build the challenge: %s", err)
	}
	if !challenge.IsAlive() {
		log.Fatal("challenge is outdated")
	}

	fmt.Printf("challenge: %v\n", []string{
		strconv.Itoa(challenge.LeadingZeroBitCount().ToInt()),
		challenge.CreatedAt().MustGet().ToString(),
		challenge.TTL().MustGet().ToString(),
		challenge.Resource().MustGet().ToString(),
		challenge.SerializedPayload().ToString(),
		challenge.Hash().Name(),
		challenge.HashDataLayout().ToString(),
	})

	solution, err := challenge.Solve(context.Background(), pow.SolveParams{
		RandomInitialNonceParams: mo.Some(powValueTypes.RandomNonceParams{
			// use `crypto/rand.Reader` in a production
			RandomReader: bytes.NewReader([]byte("dummy")),
			MinRawValue:  big.NewInt(123),
			MaxRawValue:  big.NewInt(142),
		}),
	})
	if err != nil {
		log.Fatalf("unable to solve the challenge: %s", err)
	}

	fmt.Printf("nonce: %s\n", solution.Nonce().ToString())
	fmt.Printf("hash sum: %x\n", solution.HashSum().OrEmpty().ToBytes())

	if err := solution.Verify(); err != nil {
		log.Fatalf("unable to verify the solution: %s", err)
	}

	fmt.Print("verification: OK\n")

	// Output:
	// challenge: [5 2000-01-02T03:04:05.000000006Z 876000h0m0s https://example.com/ dummy SHA-256 {{.Challenge.LeadingZeroBitCount.ToInt}}:{{.Challenge.CreatedAt.MustGet.ToString}}:{{.Challenge.TTL.MustGet.ToString}}:{{.Challenge.Resource.MustGet.ToString}}:{{.Challenge.SerializedPayload.ToString}}:{{.Challenge.Hash.Name}}:{{.Challenge.HashDataLayout.ToString}}:{{.Nonce.ToString}}]
	// nonce: 136
	// hash sum: 060056e78e0b90e48c765d4f64c0f63d5926e28a56f3cd229bdc78225f91cd51
	// verification: OK
}

func Example_withInterruption() {
	tooBigLeadingZeroBitCount, err := powValueTypes.NewLeadingZeroBitCount(100)
	if err != nil {
		log.Fatalf("unable to construct the leading zero bit count: %s", err)
	}

	namedHash, err := powValueTypes.NewHashWithName(sha256.New(), "SHA-256")
	if err != nil {
		log.Fatalf("unable to construct the hash: %s", err)
	}

	challenge, err := pow.NewChallengeBuilder().
		SetLeadingZeroBitCount(tooBigLeadingZeroBitCount).
		SetSerializedPayload(powValueTypes.NewSerializedPayload("dummy")).
		SetHash(namedHash).
		SetHashDataLayout(powValueTypes.MustParseHashDataLayout(
			"{{ .Challenge.LeadingZeroBitCount.ToInt }}" +
				":{{ .Challenge.SerializedPayload.ToString }}" +
				":{{ .Nonce.ToString }}",
		)).
		Build()
	if err != nil {
		log.Fatalf("unable to build the challenge: %s", err)
	}

	fmt.Printf("challenge: %v\n", []string{
		strconv.Itoa(challenge.LeadingZeroBitCount().ToInt()),
		challenge.SerializedPayload().ToString(),
		challenge.Hash().Name(),
		challenge.HashDataLayout().ToString(),
	})

	_, solvingErr := challenge.Solve(context.Background(), pow.SolveParams{
		MaxAttemptCount: mo.Some(1000),
	})
	if solvingErr == nil {
		log.Fatal("solving must fail")
	}
	if !errors.Is(solvingErr, powErrors.ErrTaskInterruption) {
		log.Fatalf("unexpected solving error: %s", solvingErr)
	}

	fmt.Print("solving: interrupted\n")

	// Output:
	// challenge: [100 dummy SHA-256 {{.Challenge.LeadingZeroBitCount.ToInt}}:{{.Challenge.SerializedPayload.ToString}}:{{.Nonce.ToString}}]
	// solving: interrupted
}
