// Migrated from: application.yml + SecurityConfig.java
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	ServerPort       string
	DBHost           string
	DBPort           string
	DBName           string
	DBUser           string
	DBPassword       string
	ActuatorUsername  string
	ActuatorPassword string
	AppName          string
	AppVersion       string
	BuildTime        string
	GitCommit        string
	SessionSecret    string
	CookieSecure     bool
	OIDCIssuerURL    string
	OIDCClientID     string
	OIDCClientSecret string
	OIDCRedirectURL  string
	FrontendURL      string
	AppleClientID    string
}

// sensitive keys whose values should be masked in log output
var sensitiveKeys = map[string]bool{
	"DB_PASSWORD":      true,
	"OIDC_CLIENT_SECRET": true,
	"SESSION_SECRET":   true,
	"ACTUATOR_PASSWORD": true,
}

func Load() *Config {
	entries := []configEntry{
		{"SERVER_PORT", "8080", false},
		{"DB_HOST", "localhost", false},
		{"DB_PORT", "3306", false},
		{"DB_NAME", "coffeediary", false},
		{"DB_USER", "app", false},
		{"DB_PASSWORD", "app", false},
		{"ACTUATOR_USERNAME", "actuator", false},
		{"ACTUATOR_PASSWORD", "changeme", false},
		{"APP_NAME", "coffee-diary-backend", false},
		{"APP_VERSION", "0.0.1-SNAPSHOT", false},
		{"BUILD_TIME", "unknown", false},
		{"GIT_COMMIT", "unknown", false},
		{"SESSION_SECRET", "change-me-in-production-32chars!", false},
		{"COOKIE_SECURE", "false", false},
		{"OIDC_ISSUER_URL", "https://id.oglimmer.de/realms/oglimmer", false},
		{"OIDC_CLIENT_ID", "coffee-diary", false},
		{"OIDC_CLIENT_SECRET", "", true},
		{"OIDC_REDIRECT_URL", "http://localhost:8080/api/auth/callback", false},
		{"FRONTEND_URL", "http://localhost:5173", false},
		{"APPLE_CLIENT_ID", "com.oglimmer.CoffeeDiary", false},
	}

	values := make(map[string]string, len(entries))
	var missing []string

	for _, e := range entries {
		v := os.Getenv(e.key)
		if v == "" {
			if e.required {
				missing = append(missing, e.key)
				continue
			}
			v = e.defaultVal
		}
		values[e.key] = v

		displayVal := v
		if sensitiveKeys[e.key] {
			displayVal = maskValue(v)
		}
		source := "env"
		if os.Getenv(e.key) == "" {
			source = "default"
		}
		slog.Info("config", "key", e.key, "value", displayVal, "source", source)
	}

	if len(missing) > 0 {
		slog.Error("required environment variables not set", "missing", strings.Join(missing, ", "))
		os.Exit(1)
	}

	return &Config{
		ServerPort:       values["SERVER_PORT"],
		DBHost:           values["DB_HOST"],
		DBPort:           values["DB_PORT"],
		DBName:           values["DB_NAME"],
		DBUser:           values["DB_USER"],
		DBPassword:       values["DB_PASSWORD"],
		ActuatorUsername:  values["ACTUATOR_USERNAME"],
		ActuatorPassword: values["ACTUATOR_PASSWORD"],
		AppName:          values["APP_NAME"],
		AppVersion:       values["APP_VERSION"],
		BuildTime:        values["BUILD_TIME"],
		GitCommit:        values["GIT_COMMIT"],
		SessionSecret:    values["SESSION_SECRET"],
		CookieSecure:     values["COOKIE_SECURE"] == "true",
		OIDCIssuerURL:    values["OIDC_ISSUER_URL"],
		OIDCClientID:     values["OIDC_CLIENT_ID"],
		OIDCClientSecret: values["OIDC_CLIENT_SECRET"],
		OIDCRedirectURL:  values["OIDC_REDIRECT_URL"],
		FrontendURL:      values["FRONTEND_URL"],
		AppleClientID:    values["APPLE_CLIENT_ID"],
	}
}

type configEntry struct {
	key        string
	defaultVal string
	required   bool
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func maskValue(v string) string {
	if len(v) <= 4 {
		return "****"
	}
	return v[:2] + "****" + v[len(v)-2:]
}
