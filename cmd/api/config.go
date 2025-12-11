package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ucok-man/tcsa/internal/validator"
)

func init() {
	env := os.Getenv("TCSA_ENV")
	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	godotenv.Load()
}

type Config struct {
	Port     uint   `mapstructure:"PORT" validate:"required,port"`
	Env      string `mapstructure:"ENV" validate:"required,oneof=development staging production"`
	Database struct {
		Dsn         string        `mapstructure:"DB_DSN" validate:"required,url"`
		MaxOpenConn int           `mapstructure:"DB_MAX_OPEN_CONN" validate:"required,min=1,max=100"`
		MaxIdleConn int           `mapstructure:"DB_MAX_IDLE_CONN" validate:"required,min=1,max=100"`
		MaxIdleTime time.Duration `mapstructure:"DB_MAX_IDLE_TIME" validate:"required,min=1s"`
	} `mapstructure:",squash"`
	Log struct {
		Level string `mapstructure:"LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	} `mapstructure:",squash"`
	Cors struct {
		TrustedOrigins []string `mapstructure:"CORS_TRUSTED_ORIGINS" validate:"omitempty,dive,url"`
	} `mapstructure:",squash"`
}

func NewConfig() (Config, error) {
	// Bind environment variables
	viper.SetEnvPrefix("TCSA")
	viper.AutomaticEnv()

	pflag.Uint("port", 3000, "API server port")
	pflag.String("env", "development", "Environment (development/staging/production)")
	pflag.String("db-dsn", "", "Database connection string")
	pflag.Int("db-max-open-conn", 25, "Database max open connections")
	pflag.Int("db-max-idle-conn", 25, "Database max idle connections")
	pflag.Duration("db-max-idle-time", 15*time.Minute, "Database max idle time")
	pflag.String("log-level", "debug", "Log level (debug/info/warn/error)")
	pflag.StringSlice("cors-trusted-origins", []string{}, "Trusted CORS origins (comma separated)")

	pflag.Parse()

	// Bind flags to Viper keys, flags override environment
	viper.BindPFlag("PORT", pflag.Lookup("port"))
	viper.BindPFlag("ENV", pflag.Lookup("env"))
	viper.BindPFlag("DB_DSN", pflag.Lookup("db-dsn"))
	viper.BindPFlag("DB_MAX_OPEN_CONN", pflag.Lookup("db-max-open-conn"))
	viper.BindPFlag("DB_MAX_IDLE_CONN", pflag.Lookup("db-max-idle-conn"))
	viper.BindPFlag("DB_MAX_IDLE_TIME", pflag.Lookup("db-max-idle-time"))
	viper.BindPFlag("LOG_LEVEL", pflag.Lookup("log-level"))
	viper.BindPFlag("CORS_TRUSTED_ORIGINS", pflag.Lookup("cors-trusted-origins"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unable to decode config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}
