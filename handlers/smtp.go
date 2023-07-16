package handlers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var Smtp *gomail.Dialer

func ConnectSmtp() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file!")
	}
	Smtp = gomail.NewDialer(os.Getenv("SMTP_HOST"), 587, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"))
}
