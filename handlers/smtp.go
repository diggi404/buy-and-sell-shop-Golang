package handlers

import (
	"gopkg.in/gomail.v2"
)

var Smtp *gomail.Dialer

func ConnectSmtp() {
	Smtp = gomail.NewDialer("smtp.gmail.com", 587, "karianfavreau9@gmail.com", "bcmvowqxuqomtmza")
}
