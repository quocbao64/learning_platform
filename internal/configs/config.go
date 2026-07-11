package configs

import "github.com/spf13/viper"

type Config struct {
	DbConfig        DbConfig        `mapstructure:",squash"`
	MigrationConfig MigrationConfig `mapstructure:",squash"`
	AppConfig       AppConfig       `mapstructure:",squash"`
	JWTConfig       JWTConfig       `mapstructure:",squash"`
}

type DbConfig struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

type MigrationConfig struct {
	MigrationPath string `mapstructure:"MIGRATION_PATH"`
}

type AppConfig struct {
	AppPort string `mapstructure:"APP_PORT"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"SECRET"`
	TTLMinutes int    `mapstructure:"TTL_MINUTES"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
