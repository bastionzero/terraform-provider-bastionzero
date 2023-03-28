package sweep

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bastionzero/bastionzero-sdk-go/bastionzero"
)

const TestNamePrefix = "tf-acc-test-"

// SweeperClient returns a common provider client to clean leftover BastionZero
// resources from acceptance tests that did not run correctly
func SweeperClient() (*bastionzero.Client, error) {
	opts := []bastionzero.ClientOpt{}
	apiSecret := os.Getenv("BASTIONZERO_API_SECRET")
	if apiSecret == "" {
		return nil, fmt.Errorf("empty BASTIONZERO_API_SECRET")
	}
	host := os.Getenv("BASTIONZERO_HOST")
	// If custom host specified, configure client with base URL
	if host != "" {
		opts = append(opts, bastionzero.WithBaseURL(host))
	}

	// Create a new BastionZero client using the configuration values
	client, err := bastionzero.NewFromAPISecret(http.DefaultClient, apiSecret, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating BastionZero client")
	}

	return client, nil
}
