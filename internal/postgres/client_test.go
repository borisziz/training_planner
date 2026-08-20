package postgres

import (
	"errors"
	"testing"
	"time"
)

func TestParseConfig(t *testing.T) {
	t.Run("rejects empty connection string", func(t *testing.T) {
		_, err := parseConfig("  ")
		if !errors.Is(err, ErrEmptyConnectionString) {
			t.Fatalf("expected ErrEmptyConnectionString, got %v", err)
		}
	})

	t.Run("applies safe connect timeout by default", func(t *testing.T) {
		config, err := parseConfig("postgres://user:password@localhost:5432/training")
		if err != nil {
			t.Fatalf("parse config: %v", err)
		}

		if config.ConnectTimeout != defaultConnectTimeout {
			t.Fatalf(
				"expected default timeout %s, got %s",
				defaultConnectTimeout,
				config.ConnectTimeout,
			)
		}
	})

	t.Run("preserves explicit connection settings", func(t *testing.T) {
		config, err := parseConfig(
			"postgres://user:password@localhost:5432/training" +
				"?connect_timeout=9",
		)
		if err != nil {
			t.Fatalf("parse config: %v", err)
		}

		if config.ConnectTimeout != 9*time.Second {
			t.Errorf("expected connect timeout 9s, got %s", config.ConnectTimeout)
		}
	})
}
