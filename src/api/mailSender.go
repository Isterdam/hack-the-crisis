package api

import (
	"fmt"
	"net/smtp"
	"net/mail"
	"encoding/base64"
	"strings"
	"os"
)

func encodeRFC2047(String string) string{
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(strings.Trim(addr.String(), " <>"), "<@")
}

func SendMail(to, subject, content string) {
	mail := "booklie.pass@gmail.com"
	password := os.Getenv("MAILPASS")

	header := make(map[string]string)
	header["To"] = to
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	msg := "From: Booklie <" + mail + "> \n"
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(content))

	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", mail, password, "smtp.gmail.com"), "Booklie", []string{to}, []byte(msg))

	if err != nil {
		fmt.Println(err)
	}
}
