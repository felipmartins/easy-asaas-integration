package config

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoadAppliesSafeDefaults(t *testing.T) {
	config, err := load(testEnvironment(map[string]string{
		"DATABASE_URL":        "postgres://test:test@localhost:5432/easy_asaas_test",
		"ASAAS_API_KEY":       "sandbox-test-api-key",
		"ASAAS_WEBHOOK_TOKEN": "test-webhook-token",
	}))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Environment != EnvironmentDevelopment {
		t.Errorf("Environment = %q, want %q", config.Environment, EnvironmentDevelopment)
	}
	if config.HTTPAddress != ":8080" {
		t.Errorf("HTTPAddress = %q, want %q", config.HTTPAddress, ":8080")
	}
	if config.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %s, want %s", config.ShutdownTimeout, 10*time.Second)
	}
	if config.Asaas.Environment != AsaasEnvironmentSandbox {
		t.Errorf("Asaas.Environment = %q, want %q", config.Asaas.Environment, AsaasEnvironmentSandbox)
	}
}

func TestLoadParsesExplicitValues(t *testing.T) {
	config, err := load(testEnvironment(map[string]string{
		"APP_ENV":             "production",
		"HTTP_ADDR":           "127.0.0.1:9090",
		"SHUTDOWN_TIMEOUT":    "25s",
		"DATABASE_URL":        "postgres://test:test@localhost:5432/easy_asaas_test",
		"ASAAS_ENVIRONMENT":   "production",
		"ASAAS_API_KEY":       "production-test-key",
		"ASAAS_WEBHOOK_TOKEN": "production-test-webhook-token",
	}))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Environment != EnvironmentProduction {
		t.Errorf("Environment = %q, want %q", config.Environment, EnvironmentProduction)
	}
	if config.HTTPAddress != "127.0.0.1:9090" {
		t.Errorf("HTTPAddress = %q, want %q", config.HTTPAddress, "127.0.0.1:9090")
	}
	if config.ShutdownTimeout != 25*time.Second {
		t.Errorf("ShutdownTimeout = %s, want %s", config.ShutdownTimeout, 25*time.Second)
	}
	if config.Asaas.Environment != AsaasEnvironmentProduction {
		t.Errorf("Asaas.Environment = %q, want %q", config.Asaas.Environment, AsaasEnvironmentProduction)
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	base := map[string]string{
		"DATABASE_URL":        "postgres://test:test@localhost:5432/easy_asaas_test",
		"ASAAS_API_KEY":       "sandbox-test-api-key",
		"ASAAS_WEBHOOK_TOKEN": "test-webhook-token",
	}
	tests := []struct {
		name    string
		overlay map[string]string
		wantErr string
	}{
		{
			name:    "unsupported app environment",
			overlay: map[string]string{"APP_ENV": "staging"},
			wantErr: "APP_ENV must be development, test, or production",
		},
		{
			name:    "unsupported Asaas environment",
			overlay: map[string]string{"ASAAS_ENVIRONMENT": "live"},
			wantErr: "ASAAS_ENVIRONMENT must be sandbox or production",
		},
		{
			name:    "invalid HTTP address",
			overlay: map[string]string{"HTTP_ADDR": "localhost"},
			wantErr: "invalid HTTP_ADDR",
		},
		{
			name:    "invalid port",
			overlay: map[string]string{"HTTP_ADDR": ":70000"},
			wantErr: "port must be between 1 and 65535",
		},
		{
			name:    "invalid shutdown timeout",
			overlay: map[string]string{"SHUTDOWN_TIMEOUT": "soon"},
			wantErr: "invalid SHUTDOWN_TIMEOUT",
		},
		{
			name:    "non-positive shutdown timeout",
			overlay: map[string]string{"SHUTDOWN_TIMEOUT": "0s"},
			wantErr: "SHUTDOWN_TIMEOUT must be greater than zero",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := cloneValues(base)
			for name, value := range test.overlay {
				values[name] = value
			}

			_, err := load(testEnvironment(values))
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("load error = %v, want error containing %q", err, test.wantErr)
			}
		})
	}
}

func TestLoadRequiresSecretsWithoutIncludingValuesInErrors(t *testing.T) {
	config, err := load(testEnvironment(map[string]string{}))
	if err == nil {
		t.Fatal("load config succeeded without required values")
	}
	if config != (Config{}) {
		t.Fatalf("invalid configuration returned partial values: %#v", config)
	}

	for _, name := range []string{"DATABASE_URL", "ASAAS_API_KEY", "ASAAS_WEBHOOK_TOKEN"} {
		if !strings.Contains(err.Error(), name+" is required") {
			t.Errorf("error %q does not report missing %s", err, name)
		}
	}
}

func TestSecretFormattingRedactsValue(t *testing.T) {
	const sensitiveValue = "synthetic-secret-value"
	secret := Secret{value: sensitiveValue}
	config := Config{
		Environment:     EnvironmentTest,
		HTTPAddress:     ":8080",
		ShutdownTimeout: time.Second,
		DatabaseURL:     Secret{value: "postgres://synthetic-user:synthetic-password@localhost/test"},
		Asaas: AsaasConfig{
			Environment:  AsaasEnvironmentSandbox,
			APIKey:       secret,
			WebhookToken: Secret{value: "synthetic-webhook-secret"},
		},
	}

	for _, formatted := range []string{
		fmt.Sprintf("%v", secret),
		fmt.Sprintf("%#v", secret),
		fmt.Sprintf("%v", config),
		fmt.Sprintf("%+v", config),
		fmt.Sprintf("%#v", config),
	} {
		for _, forbidden := range []string{
			sensitiveValue,
			"synthetic-password",
			"synthetic-webhook-secret",
		} {
			if strings.Contains(formatted, forbidden) {
				t.Errorf("formatted value leaked %q: %s", forbidden, formatted)
			}
		}
	}

	if secret.Reveal() != sensitiveValue {
		t.Fatal("Reveal did not return the original secret")
	}
}

func testEnvironment(values map[string]string) lookupEnv {
	return func(name string) (string, bool) {
		value, ok := values[name]
		return value, ok
	}
}

func cloneValues(values map[string]string) map[string]string {
	cloned := make(map[string]string, len(values))
	for name, value := range values {
		cloned[name] = value
	}
	return cloned
}
