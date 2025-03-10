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

	// Loop untuk bind environment variables
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

	// Menampilkan konfigurasi yang telah dimuat dengan lebih rinci
	fmt.Println("=== Loaded Configuration ===")
	// Cetak AppConfig
	fmt.Printf("App Name: %s\n", Cfg.App.Name)
	fmt.Printf("App Port: %s\n", Cfg.App.Port)
	fmt.Printf("Encryption Salt: %d\n", Cfg.App.Encryption.Salt)
	fmt.Printf("JWT Secret: %s\n", Cfg.App.Encryption.JWTSecret)

	// Cetak DBConfig
	fmt.Printf("DB Host: %s\n", Cfg.DB.Host)
	fmt.Printf("DB Port: %s\n", Cfg.DB.Port)
	fmt.Printf("DB Name: %s\n", Cfg.DB.Name)
	fmt.Printf("DB User: %s\n", Cfg.DB.User)
	fmt.Printf("DB Password: %s\n", Cfg.DB.Password)
	fmt.Printf("DB Connection Pool: %+v\n", Cfg.DB.ConnectionPool)

	// Cetak connection pool details
	fmt.Printf("Max Idle Connection: %d\n", Cfg.DB.ConnectionPool.MaxIdleConnection)
	fmt.Printf("Max Open Connection: %d\n", Cfg.DB.ConnectionPool.MaxOpenConnection)
	fmt.Printf("Max Lifetime Connection: %d\n", Cfg.DB.ConnectionPool.MaxLifetimeConnection)
	fmt.Printf("Max Idle Time Connection: %d\n", Cfg.DB.ConnectionPool.MaxIdletimeConnection)

	return nil
}
