# loops-go-sdk

[![CI](https://github.com/Whats-A-MattR/loops-go-sdk/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Whats-A-MattR/loops-go-sdk/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Whats-A-MattR/loops-go-sdk)](https://goreportcard.com/report/github.com/Whats-A-MattR/loops-go-sdk)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/Whats-A-MattR/loops-go-sdk)](https://pkg.go.dev/github.com/Whats-A-MattR/loops-go-sdk)

Community-maintained Go SDK for the [Loops](https://loops.so) API. Server-side only; types and behaviour follow the [Loops OpenAPI spec](https://app.loops.so/openapi.json).

## Install

```bash
go get github.com/Whats-A-MattR/loops-go-sdk
```

Requires Go 1.21+.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Whats-A-MattR/loops-go-sdk"
)

func main() {
	client := loops.NewClient("your-api-key")
	ctx := context.Background()

	// Verify API key and get team name
	resp, err := client.GetAPIKey(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Team:", resp.TeamName)
}
```

## Examples

### Create a contact

```go
client := loops.NewClient(os.Getenv("LOOPS_API_KEY"))
ctx := context.Background()

created, err := client.CreateContact(ctx, &loops.ContactRequest{
	Email:     "user@example.com",
	FirstName: "Jane",
	LastName:  "Doe",
})
if err != nil {
	log.Fatal(err)
}
fmt.Println("Contact ID:", created.ID)
```

### Find a contact

```go
contacts, err := client.FindContact(ctx, "user@example.com", "")
if err != nil {
	log.Fatal(err)
}
for _, c := range contacts {
	fmt.Printf("%s: %s\n", c.ID, c.Email)
}
```

### Send an event (trigger emails)

```go
_, err := client.SendEvent(ctx, &loops.EventRequest{
	Email:      "user@example.com",
	EventName:  "signed_up",
	EventProperties: map[string]interface{}{
		"plan": "pro",
	},
}, "") // optional idempotency key as second argument
if err != nil {
	log.Fatal(err)
}
```

### Send a transactional email

```go
_, err := client.SendTransactional(ctx, &loops.TransactionalRequest{
	Email:           "user@example.com",
	TransactionalID: "clxxxxxxxxxxxx", // ID from your Loops dashboard
	DataVariables: map[string]interface{}{
		"name": "Jane",
		"resetLink": "https://example.com/reset",
	},
}, "")
if err != nil {
	log.Fatal(err)
}
```

### Handle API errors

```go
_, err := client.CreateContact(ctx, &loops.ContactRequest{Email: "user@example.com"})
if err != nil {
	var apiErr *loops.APIError
	if errors.As(err, &apiErr) {
		fmt.Println("Status:", apiErr.StatusCode, "Message:", apiErr.Message)
		return
	}
	log.Fatal(err)
}
```
Use the standard `errors` package for `errors.As`.

### Custom base URL or HTTP client

```go
client := loops.NewClient(apiKey,
	loops.WithBaseURL("https://custom.example.com/api/v1"),
	loops.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

## API overview

| Area | Methods |
|------|--------|
| **API key** | `GetAPIKey` |
| **Contacts** | `CreateContact`, `UpdateContact`, `FindContact`, `DeleteContact` |
| **Contact properties** | `CreateContactProperty`, `ListContactProperties` |
| **Mailing lists** | `GetLists` |
| **Events** | `SendEvent` |
| **Transactional** | `SendTransactional`, `ListTransactionals` |
| **Dedicated IPs** | `GetDedicatedSendingIPs` |

All methods take a `context.Context` as the first argument. See [pkg.go.dev](https://pkg.go.dev/github.com/Whats-A-MattR/loops-go-sdk) for full API docs.

## GitHub

- **CI** runs on every push and pull request (build, tests, fuzzing, vet).
- **OpenAPI spec check** runs daily; if the [live Loops spec](https://app.loops.so/openapi.json) changes, the workflow fails and opens an issue so the SDK can be updated.
- **Branch protection:** See [.github/BRANCH_PROTECTION.md](.github/BRANCH_PROTECTION.md) for recommended settings.

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on opening issues, submitting PRs, running tests, and keeping the SDK in sync with the Loops API.

## License

Unlicense. See [LICENSE](LICENSE).
