package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"db"`
}

type AppConfig struct {
	Name       string           `mapstructure:"name"`
	Port       string           `mapstructure:"port"`
	Encryption EncryptionConfig `mapstructure:"encryption"`
}

type EncryptionConfig struct {
	Salt      uint8  `mapstructure:"salt"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

type DBConfig struct {
	Host           string                 `mapstructure:"host"`
	Port           string                 `mapstructure:"port"`
	User           string                 `mapstructure:"user"`
	Password       string                 `mapstructure:"password"`
	Name           string                 `mapstructure:"name"`
	ConnectionPool DBConnectionPoolConfig `mapstructure:"connection_pool"`
}

type DBConnectionPoolConfig struct {
	MaxIdleConnection     uint8 `mapstructure:"max_idle_connection"`
	MaxOpenConnection     uint8 `mapstructure:"max_open_connection"`
	MaxLifetimeConnection uint8 `mapstructure:"max_lifetime_connection"`
	MaxIdletimeConnection uint8 `mapstructure:"max_idletime_connection"`
}

var Cfg Config

func LoadConfig(filename string) error {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	// Bind environment variables agar bisa digunakan
	viper.AutomaticEnv()

	// Secara eksplisit bind environment variable
	envVars := map[string]string{
		"app.encryption.jwt_secret": "JWT_SECRET",
		"db.host":                   "PGHOST",
		"db.port":                   "PGPORT",
		"db.user":                   "PGUSER",
		"db.password":               "PGPASSWORD",
		"db.name":                   "PGDATABASE",
	}

	for key, envVar := range envVars {
		if err := viper.BindEnv(key, envVar); err != nil {
			return fmt.Errorf("error binding env %s: %w", envVar, err)
		}
	}

	// Baca file konfigurasi
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal ke struct Config
	if err := viper.Unmarshal(&Cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	fmt.Printf("Config loaded: %+v\n", Cfg) // Debugging log
	return nil
}
