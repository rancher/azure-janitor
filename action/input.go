package action

import (
	"fmt"

	"github.com/caarlos0/env/v9"
	"go.uber.org/multierr"
)

type Input struct {
	ResourceGroupPattern string `env:"INPUT_RESOURCE_GROUPS"`
	AzureSubscriptionID  string `env:"INPUT_SUBSCRIPTION_ID"`
	AzureClientID        string `env:"INPUT_CLIENT_ID"`
	AzureClientSecret    string `env:"INPUT_CLIENT_SECRET"`
	AzureTenantID        string `env:"INPUT_TENANT_ID"`
	Commit               bool   `env:"INPUT_COMMIT"`
}

// NewInput creates a new input from the environment variables.
func NewInput() (*Input, error) {
	input := &Input{}
	if err := env.Parse(input); err != nil {
		return nil, fmt.Errorf("parsing environment variables: %w", err)
	}

	return input, nil
}

func (i *Input) Validate() error {
	var err error

	if i.AzureClientID == "" {
		err = multierr.Append(err, ErrClientIDRequired)
	}
	if i.AzureClientSecret == "" {
		err = multierr.Append(err, ErrClientSecretRequired)
	}
	if i.AzureSubscriptionID == "" {
		err = multierr.Append(err, ErrSubscriptionIDRequired)
	}
	if i.AzureTenantID == "" {
		err = multierr.Append(err, ErrTenantIDRequired)
	}
	if i.ResourceGroupPattern == "" {
		err = multierr.Append(err, ErrResourceGrpPatternRequired)
	}

	return err
}
