package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	AppConfig        AppConfig
	DBConfig         DBConfig
	RedisConfig      RedisConfig
	JWTConfig        JWTConfig
	MonitoringConfig MonitoringConfig
	MailConfig       MailConfig
}

type AppConfig struct {
	Name            string
	Version         string
	Port            int
	Env             string
	ShutdownTimeout int
	LogDir          string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host          string
	Port          int
	Password      string
	DB            int
	PoolSize      int
	MinIdleConns  int
	MaxRetries    int
	RetryInterval int
}

type JWTConfig struct {
	Secret            string
	RefreshSecret     string
	Expiration        int
	RefreshExpiration int
}

type MonitoringConfig struct {
	Prometheus PrometheusConfig
	Grafana    Grafana
}

type PrometheusConfig struct {
	Enabled        bool
	Endpoint       string
	Port           int
	ScrapeInterval string
	RetentionTime  string
}

type Grafana struct {
	Enabled          bool
	Port             int
	AdminUser        string
	AdminPassword    string
	DashboardPath    string
	ProvisioningPath string
}

type MailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPPassword string
	FromEmail    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &Config{
		AppConfig: AppConfig{
			Name:            getEnv("APP_NAME", "motorbike-rental-backend-v2"),
			Env:             getEnv("APP_ENV", "development"),
			Port:            getEnvAsInt("APP_PORT", 3005),
			Version:         getEnv("APP_VERSION", "1.0.0"),
			ShutdownTimeout: getEnvAsInt("APP_SHUTDOWN_TIMEOUT", 5),
			LogDir:          getEnv("APP_LOG_DIR", "./logs"),
		},
		DBConfig: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "motorbike-rental-backend-v2"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		RedisConfig: RedisConfig{
			Host:          getEnv("REDIS_HOST", "localhost"),
			Port:          getEnvAsInt("REDIS_PORT", 6379),
			Password:      getEnv("REDIS_PASSWORD", ""),
			DB:            getEnvAsInt("REDIS_DB", 0),
			PoolSize:      getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns:  getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:    getEnvAsInt("REDIS_MAX_RETRIES", 3),
			RetryInterval: getEnvAsInt("REDIS_RETRY_INTERVAL", 100),
		},
		JWTConfig: JWTConfig{
			Secret:            getEnv("JWT_SECRET", ""),
			RefreshSecret:     getEnv("JWT_REFRESH_SECRET", ""),
			Expiration:        getEnvAsInt("JWT_EXPIRATION", 15),
			RefreshExpiration: getEnvAsInt("JWT_REFRESH_EXPIRATION", 30),
		},
		MonitoringConfig: MonitoringConfig{
			Prometheus: PrometheusConfig{
				Enabled:        getEnvAsBool("PROMETHEUS_ENABLED", true),
				Endpoint:       getEnv("PROMETHEUS_ENDPOINT", "/metrics"),
				Port:           getEnvAsInt("PROMETHEUS_PORT", 9090),
				ScrapeInterval: getEnv("PROMETHEUS_SCRAPE_INTERVAL", "15s"),
				RetentionTime:  getEnv("PROMETHEUS_RETENTION_TIME", "15d"),
			},
			Grafana: Grafana{
				Enabled:          getEnvAsBool("GRAFANA_ENABLED", true),
				Port:             getEnvAsInt("GRAFANA_PORT", 3000),
				AdminUser:        getEnv("GRAFANA_ADMIN_USER", "admin"),
				AdminPassword:    getEnv("GRAFANA_ADMIN_PASSWORD", "admin"),
				DashboardPath:    getEnv("GRAFANA_DASHBOARD_PATH", "./grafana/dashboards"),
				ProvisioningPath: getEnv("GRAFANA_PROVISIONING_PATH", "./grafana/provisioning"),
			},
		},
		MailConfig: MailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("SMTP_FROM_EMAIL", ""),
		},
	}

	return config, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
