package action

import "errors"

var (
	ErrResourceGrpPatternRequired = errors.New("resource group pattern is required")
	ErrSubscriptionIDRequired     = errors.New("azure subscription id is required")
	ErrClientIDRequired           = errors.New("azure client id is required")
	ErrClientSecretRequired       = errors.New("azure client secret is required")
	ErrTenantIDRequired           = errors.New("azure tenant id is required")
)
