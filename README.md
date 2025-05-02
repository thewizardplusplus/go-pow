# go-pow

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-pow?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-pow)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-pow)](https://goreportcard.com/report/github.com/thewizardplusplus/go-pow)
[![lint](https://github.com/thewizardplusplus/go-pow/actions/workflows/lint.yaml/badge.svg)](https://github.com/thewizardplusplus/go-pow/actions/workflows/lint.yaml)
[![test](https://github.com/thewizardplusplus/go-pow/actions/workflows/test.yaml/badge.svg)](https://github.com/thewizardplusplus/go-pow/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-pow/graph/badge.svg?token=m3HjBxbUlg)](https://codecov.io/gh/thewizardplusplus/go-pow)

A library that implements a [proof-of-work](https://en.wikipedia.org/wiki/Proof_of_work) system with customizable challenges.

## Features

- use of patterns:
  - implementation based on [Domain-Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design) principles;
  - usage of the [Builder pattern](https://en.wikipedia.org/wiki/Builder_pattern) to ensure entity invariants;
  - definition of entities with explicitly declared [value types](https://en.wikipedia.org/wiki/Value_object) for each field;
- definition of challenges with the following fields:
  - `leading zero bit count` &mdash; the required number of leading zero bits in the resulting hash;
  - `target bit index` &mdash; a bit index used in determining the solution’s difficulty:
    - it refers to the bit position such that the resulting hash will be less than the number where that bit is set (for example, with a 256-bit hash and a requirement of 6 leading zeros, we would set the 250th bit);
    - only one of `leading zero bit count` or `target bit index` should be set explicitly &mdash; the other is derived from it;
  - `created at` _(optional)_ &mdash; the timestamp when the challenge was created;
  - `TTL` _(optional)_ &mdash; the duration after which the challenge expires:
    - `created at` and `TTL` must either be both specified or both omitted;
  - `resource` _(optional)_ &mdash; the resource associated with the challenge, typically for scoping:
    - based on the [`net/url.URL`](https://pkg.go.dev/net/url@go1.23.0#URL) type, but any [Uniform Resource Identifier (URI)](https://en.wikipedia.org/wiki/Uniform_Resource_Identifier) format is allowed;
  - `serialized payload` &mdash; the raw data to be included in the hash:
    - must be pre-serialized to a string; the library does not handle serialization itself;
  - `hash` &mdash; the hash function used to verify the solution:
    - based on the [`hash.Hash`](https://pkg.go.dev/hash@go1.23.0#Hash) interface;
  - `hash data layout` &mdash; the structure of the data used during hashing:
    - based on the [`text/template.Template`](https://pkg.go.dev/text/template@go1.23.0#Template) type;
    - defines which fields of the challenge will be hashed and in what order, giving full control over the hash input structure;
- generation of solutions that meet specified challenge criteria:
  - starting nonce value:
    - it can be zero;
    - it can be randomly selected within a given range;
  - the generation process can be interrupted via:
    - [context](https://pkg.go.dev/context@go1.23.0#Context) cancellation;
    - an attempt limit;
- validation of solutions against their corresponding challenges;
- sentinel errors provided through a dedicated `errors` subpackage.

## Installation

```
$ go get github.com/thewizardplusplus/go-pow
```

## Examples

Minimal (also see in the playground: https://go.dev/play/p/wsUGURDKFvb):

```go
package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"

	pow "github.com/thewizardplusplus/go-pow"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func main() {
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
```

Full (also see in the playground: https://go.dev/play/p/_3kPX0VtHFA):

```go
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/samber/mo"
	pow "github.com/thewizardplusplus/go-pow"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func main() {
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
```

With interruption (also see in the playground: https://go.dev/play/p/xT2G05H8VqN):

```go
package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/samber/mo"
	pow "github.com/thewizardplusplus/go-pow"
	powErrors "github.com/thewizardplusplus/go-pow/errors"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func main() {
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
```

## License

The MIT License (MIT)

Copyright &copy; 2025 thewizardplusplus
