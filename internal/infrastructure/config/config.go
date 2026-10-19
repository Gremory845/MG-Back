package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string

	JWTSecret       []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	// Correo (recuperación de contraseña). Si SMTPHost está vacío, el código se imprime en la consola.
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	MailFrom     string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no se encontró .env, usando variables de entorno del sistema")
	}

	accessMinutes := getEnvInt("ACCESS_TOKEN_TTL_MINUTES", 15)
	refreshDays := getEnvInt("REFRESH_TOKEN_TTL_DAYS", 7)

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Println("ADVERTENCIA: JWT_SECRET no está configurado, usando uno de desarrollo inseguro")
		jwtSecret = "dev-secret-cambia-esto-en-produccion"
	}

	return Config{
		DBUser:     getEnv("DB_USER", "musical_game"),
		DBPassword: getEnv("DB_PASSWORD", "devpassword"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "musical_game"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		JWTSecret:       []byte(jwtSecret),
		AccessTokenTTL:  time.Duration(accessMinutes) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshDays) * 24 * time.Hour,

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		MailFrom:     getEnv("MAIL_FROM", "Musical Game <no-reply@musicalgame.local>"),
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
