package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServicePort   int
	DbHost        string
	DbPort        int
	DbName        string
	DbUser        string
	DbPassword    string
	DbSsl         string
	JwtSecret     string
	RefreshSecret string
}

func (c *Config) LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	c.ServicePort, err = strconv.Atoi(os.Getenv("SERVICE_PORT"))
	c.DbHost = os.Getenv("DB_HOST")
	c.DbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
	c.DbName = os.Getenv("DB_NAME")
	c.DbUser = os.Getenv("DB_USER")
	c.DbPassword = os.Getenv("DB_PASSWORD")
	c.DbSsl = os.Getenv("DB_SSL")
	c.JwtSecret = os.Getenv("JWT_SECRET")
	c.RefreshSecret = os.Getenv("REFRESH_SECRET")
}
