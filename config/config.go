package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	App      AppConfig
	Postgres PostgresConfig
	Server   ServerConfig
	S3       S3Config
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PostgresAdd        string
}

type AppConfig struct {
	JwtSecret  string
	BcryptSalt string
}

type S3Config struct {
	ID         string
	SecretKey  string
	BucketName string
}

func NewConfig() *Configuration {
	additional := ""
	if os.Getenv("ENV") != "production" {
		additional = "&sslrootcert=ap-southeast-1-bundle.pem&Timezone=UTC"
		if godotenv.Load() != nil {
			fmt.Println("error loading .env file")
		}
	}

	config := Configuration{
		Server: ServerConfig{
			Port: ":8080",
		},
		Postgres: PostgresConfig{
			PostgresqlHost:     os.Getenv("DB_HOST"),
			PostgresqlPort:     os.Getenv("DB_PORT"),
			PostgresqlUser:     os.Getenv("DB_USERNAME"),
			PostgresqlDbname:   os.Getenv("DB_NAME"),
			PostgresqlPassword: os.Getenv("DB_PASSWORD"),
			PostgresqlSSLMode:  false,
			PostgresAdd:        additional,
		},
		App: AppConfig{
			JwtSecret:  os.Getenv("JWT_SECRET"),
			BcryptSalt: os.Getenv("BCRYPT_SALT"),
		},
		S3: S3Config{
			ID:         os.Getenv("S3_ID"),
			SecretKey:  os.Getenv("S3_SECRET_KEY"),
			BucketName: os.Getenv("S3_BUCKET_NAME"),
		},
	}

	return &config
}
