package config

import (
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	params := []string{
		"host=" + Get("DB_HOST", "127.0.0.1"),
		"user=" + Get("DB_USER", "postgres"),
		"port=" + Get("DB_PORT", "5432"),
		"sslmode=" + Get("DB_SSLMODE", "disable"),
		"TimeZone=" + Get("DB_TIMEZONE", "UTC"),
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		params = append(params, "password="+pass)
	}
	params = append(params, "dbname="+Get("DB_NAME", "ukm_hub_db"))
	dsn := strings.Join(params, " ")

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
