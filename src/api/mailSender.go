package api

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(to, subject, content string) {
	from := "team@shopalone.se"
	password := os.Getenv("MAILPASS")
	msg := "From: Team ShopAlone <" + from + "> \n" +
		"To: " + to + "\n" +
		"Subject:" + subject + "\n" +
		content

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Println(err)
	}
}
