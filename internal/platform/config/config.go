// Package config loads and validates process configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Environment identifies the runtime environment for this application.
type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentTest        Environment = "test"
	EnvironmentProduction  Environment = "production"
)

// AsaasEnvironment identifies which Asaas environment receives API requests.
type AsaasEnvironment string

const (
	AsaasEnvironmentSandbox    AsaasEnvironment = "sandbox"
	AsaasEnvironmentProduction AsaasEnvironment = "production"
)

// Config contains validated settings needed to start the application.
// Secret fields redact their contents when formatted for logs.
type Config struct {
	Environment     Environment
	HTTPAddress     string
	ShutdownTimeout time.Duration
	DatabaseURL     Secret
	Asaas           AsaasConfig
}

// AsaasConfig contains credentials and environment selection for Asaas.
type AsaasConfig struct {
	Environment  AsaasEnvironment
	APIKey       Secret
	WebhookToken Secret
}

// Secret stores sensitive configuration and redacts it from normal formatting.
// Call Reveal only at the boundary that needs the underlying credential.
type Secret struct {
	value string
}

// Reveal returns the secret value for use by an adapter that needs the credential.
func (s Secret) Reveal() string {
	return s.value
}

// String returns a redacted representation of the secret.
func (Secret) String() string {
	return "[REDACTED]"
}

// GoString returns a redacted Go-syntax representation of the secret.
func (Secret) GoString() string {
	return "Secret([REDACTED])"
}

// String returns a safe summary of the configuration without credential values.
func (c Config) String() string {
	return fmt.Sprintf(
		"Config{Environment:%q HTTPAddress:%q ShutdownTimeout:%s DatabaseURL:[REDACTED] Asaas:{Environment:%q APIKey:[REDACTED] WebhookToken:[REDACTED]}}",
		c.Environment,
		c.HTTPAddress,
		c.ShutdownTimeout,
		c.Asaas.Environment,
	)
}

// GoString returns a safe Go-syntax summary without credential values.
func (c Config) GoString() string {
	return c.String()
}

// Load reads process environment variables and validates the application configuration.
func Load() (Config, error) {
	return load(os.LookupEnv)
}

type lookupEnv func(string) (string, bool)

func load(lookup lookupEnv) (Config, error) {
	var validationErrors []error

	environment := Environment(defaultValue(lookup, "APP_ENV", string(EnvironmentDevelopment)))
	if !isEnvironment(environment) {
		validationErrors = append(validationErrors, errors.New("APP_ENV must be development, test, or production"))
	}

	asaasEnvironment := AsaasEnvironment(defaultValue(lookup, "ASAAS_ENVIRONMENT", string(AsaasEnvironmentSandbox)))
	if asaasEnvironment != AsaasEnvironmentSandbox && asaasEnvironment != AsaasEnvironmentProduction {
		validationErrors = append(validationErrors, errors.New("ASAAS_ENVIRONMENT must be sandbox or production"))
	}

	httpAddress := defaultValue(lookup, "HTTP_ADDR", ":8080")
	if err := validateHTTPAddress(httpAddress); err != nil {
		validationErrors = append(validationErrors, fmt.Errorf("invalid HTTP_ADDR: %w", err))
	}

	shutdownTimeout := 10 * time.Second
	if raw, ok := lookup("SHUTDOWN_TIMEOUT"); ok && strings.TrimSpace(raw) != "" {
		parsed, err := time.ParseDuration(strings.TrimSpace(raw))
		if err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err))
		} else if parsed <= 0 {
			validationErrors = append(validationErrors, errors.New("SHUTDOWN_TIMEOUT must be greater than zero"))
		} else {
			shutdownTimeout = parsed
		}
	}

	databaseURL, err := requiredSecret(lookup, "DATABASE_URL")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	asaasAPIKey, err := requiredSecret(lookup, "ASAAS_API_KEY")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	asaasWebhookToken, err := requiredSecret(lookup, "ASAAS_WEBHOOK_TOKEN")
	if err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := errors.Join(validationErrors...); err != nil {
		return Config{}, err
	}

	return Config{
		Environment:     environment,
		HTTPAddress:     httpAddress,
		ShutdownTimeout: shutdownTimeout,
		DatabaseURL:     databaseURL,
		Asaas: AsaasConfig{
			Environment:  asaasEnvironment,
			APIKey:       asaasAPIKey,
			WebhookToken: asaasWebhookToken,
		},
	}, nil
}

func isEnvironment(environment Environment) bool {
	switch environment {
	case EnvironmentDevelopment, EnvironmentTest, EnvironmentProduction:
		return true
	default:
		return false
	}
}

func validateHTTPAddress(address string) error {
	_, rawPort, err := net.SplitHostPort(address)
	if err != nil {
		return errors.New("must be a host:port address")
	}

	port, err := strconv.Atoi(rawPort)
	if err != nil || port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	return nil
}

func requiredSecret(lookup lookupEnv, name string) (Secret, error) {
	value, ok := lookup(name)
	if !ok || strings.TrimSpace(value) == "" {
		return Secret{}, fmt.Errorf("%s is required", name)
	}
	return Secret{value: strings.TrimSpace(value)}, nil
}

func defaultValue(lookup lookupEnv, name, fallback string) string {
	value, ok := lookup(name)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
