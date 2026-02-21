# regulayer-sdk-go

Official Go SDK for Regulayer. Record provable AI decisions with tamper-detectable audit trails.

## Installation

```bash
go get github.com/regulayer/regulayer-sdk-go
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/regulayer/regulayer-sdk-go"
)

func main() {
	client, err := regulayer.NewClient(regulayer.Config{
		APIKey: "rl_live_your_key_here",
	})
	if err != nil {
		panic(err)
	}

	err = client.RecordDecision(regulayer.Decision{
		System: "loan-approval-ai",
		Input: map[string]interface{}{"applicant_id": "12345", "credit_score": 720},
		Output: map[string]interface{}{"approved": true, "risk": "low"},
	})

	if err != nil {
		fmt.Printf("Failed to record decision: %v\n", err)
	} else {
		fmt.Println("Decision recorded successfully!")
	}
}
```

See [docs.regulayer.tech](https://docs.regulayer.tech) for full documentation.
