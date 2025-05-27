package config

import "errors"

// Config Errors.
var (
	ErrMissingNetBirdToken = errors.New("NetBird token is missing from provider configuration")
	ErrMissingNetBirdURL   = errors.New("NetBird URL is missing from provider configuration")
	ErrNilProviderConfig   = errors.New("provider configuration is nil")
)

// Resource Errors.
var (
	ErrGetProviderURL = errors.New("error getting provider URL")
)
