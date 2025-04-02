package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost         string
	Port               string
	DBUser             string
	DBPassword         string
	DBAddress          string
	DBName             string
	SecretJWT          string
	DBPort             string
	DBHost             string
	GoogleCliendID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func LoadEnv() Config {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory", err)
	}

	fmt.Println("current directory is:", cwd)

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: No .env file found ")
	}

	return Config{
		PublicHost:         os.Getenv("PUBLIC_HOST"),
		Port:               os.Getenv("PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		DBAddress:          fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		SecretJWT:          os.Getenv("SECRET"),
		DBPort:             os.Getenv("DB_PORT"),
		DBHost:             os.Getenv("DB_HOST"),
		GoogleCliendID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}

}

var ENV = LoadEnv()
