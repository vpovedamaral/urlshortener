package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config représente la configuration complète de l'application.
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Analytics AnalyticsConfig `mapstructure:"analytics"`
	Monitor   MonitorConfig   `mapstructure:"monitor"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	Name string `mapstructure:"name"`
}

type AnalyticsConfig struct {
	BufferSize  int `mapstructure:"buffer_size"`
	WorkerCount int `mapstructure:"worker_count"`
}

type MonitorConfig struct {
	IntervalMinutes int `mapstructure:"interval_minutes"`
}

// LoadConfig charge la configuration depuis le fichier config.yaml ou utilise les valeurs par défaut.
func LoadConfig() (*Config, error) {
	// Configure le chemin et le nom du fichier de configuration.
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Définit les valeurs par défaut.
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("database.name", "url_shortener_grp5.db")
	viper.SetDefault("analytics.buffer_size", 1000)
	viper.SetDefault("analytics.worker_count", 5)
	viper.SetDefault("monitor.interval_minutes", 5)

	// Lit le fichier de configuration.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Fichier de configuration non trouvé, utilisation des valeurs par défaut")
		} else {
			return nil, fmt.Errorf("erreur lors de la lecture du fichier de configuration: %w", err)
		}
	}

	// Mappe la configuration dans la structure Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erreur lors du demap de la configuration: %w", err)
	}

	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil
}
