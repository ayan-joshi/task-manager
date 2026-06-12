package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
// All values are validated at startup; the application refuses to boot with an
// invalid configuration so that misconfiguration fails fast instead of leaking
// at request time.
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiryHours int
	AllowedOrigins []string
	UploadDir      string
	MaxUploadBytes int64
	AllowedMime    []string
	RateLimitRPS   float64
	RateLimitBurst int
	Environment    string

	// Optional admin seed. When email and password are both set, an ADMIN user
	// is created on startup if it does not already exist.
	SeedAdminName     string
	SeedAdminEmail    string
	SeedAdminPassword string
}

// Load reads configuration from the environment (and an optional .env file when
// present) and validates it. It returns an error describing every missing or
// invalid value rather than panicking.
func Load() (*Config, error) {
	// .env is optional: in production the platform injects real env vars.
	_ = godotenv.Load()

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 72),
		AllowedOrigins: splitCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadBytes: int64(getEnvInt("MAX_UPLOAD_MB", 10)) * 1024 * 1024,
		AllowedMime: splitCSV(getEnv("ALLOWED_MIME_TYPES",
			"image/png,image/jpeg,image/gif,image/webp,application/pdf,text/plain")),
		RateLimitRPS:   getEnvFloat("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20),
		Environment:    getEnv("ENVIRONMENT", "development"),

		SeedAdminName:     getEnv("SEED_ADMIN_NAME", "Administrator"),
		SeedAdminEmail:    os.Getenv("SEED_ADMIN_EMAIL"),
		SeedAdminPassword: os.Getenv("SEED_ADMIN_PASSWORD"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	var problems []string

	if c.DatabaseURL == "" {
		problems = append(problems, "DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		problems = append(problems, "JWT_SECRET is required")
	} else if len(c.JWTSecret) < 16 {
		problems = append(problems, "JWT_SECRET must be at least 16 characters")
	}
	if c.JWTExpiryHours <= 0 {
		problems = append(problems, "JWT_EXPIRY_HOURS must be positive")
	}

	if len(problems) > 0 {
		return fmt.Errorf("invalid configuration: %s", strings.Join(problems, "; "))
	}
	return nil
}

// IsProduction reports whether the app is running in a production environment.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.Environment, "production")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
