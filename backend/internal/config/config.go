// Migrated from: application.yml + SecurityConfig.java
package config

import (
	"fmt"
	"os"
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

func Load() *Config {
	return &Config{
		ServerPort:       envOrDefault("SERVER_PORT", "8080"),
		DBHost:           envOrDefault("DB_HOST", "localhost"),
		DBPort:           envOrDefault("DB_PORT", "3306"),
		DBName:           envOrDefault("DB_NAME", "coffeediary"),
		DBUser:           envOrDefault("DB_USER", "app"),
		DBPassword:       envOrDefault("DB_PASSWORD", "app"),
		ActuatorUsername:  envOrDefault("ACTUATOR_USERNAME", "actuator"),
		ActuatorPassword:  envOrDefault("ACTUATOR_PASSWORD", "changeme"),
		AppName:          envOrDefault("APP_NAME", "coffee-diary-backend"),
		AppVersion:       envOrDefault("APP_VERSION", "0.0.1-SNAPSHOT"),
		BuildTime:        envOrDefault("BUILD_TIME", "unknown"),
		GitCommit:        envOrDefault("GIT_COMMIT", "unknown"),
		SessionSecret:    envOrDefault("SESSION_SECRET", "change-me-in-production-32chars!"),
		CookieSecure:     envOrDefault("COOKIE_SECURE", "false") == "true",
		OIDCIssuerURL:    envOrDefault("OIDC_ISSUER_URL", "https://id.oglimmer.de/realms/oglimmer"),
		OIDCClientID:     envOrDefault("OIDC_CLIENT_ID", "coffee-diary"),
		OIDCClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		OIDCRedirectURL:  envOrDefault("OIDC_REDIRECT_URL", "http://localhost:8080/api/auth/callback"),
		FrontendURL:      envOrDefault("FRONTEND_URL", "http://localhost:5173"),
		AppleClientID:    envOrDefault("APPLE_CLIENT_ID", "com.oglimmer.CoffeeDiary"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
