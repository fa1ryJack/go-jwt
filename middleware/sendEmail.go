package middleware

import (
	"log"
	"os"

	"net/smtp"
)

func SendEmail(email string) {

// Choose auth method and set it up

auth := smtp.PlainAuth("", os.Getenv("GMAIL"), os.Getenv("GMAILAPPPASSWORD"), "smtp.gmail.com")

// Here we do it all: connect to our server, set up a message and send it

to := []string{email}

msg := []byte("To: " + email + "\r\n" +

"Subject: email warning: Your IP has changed\r\n" +

"\r\n" +

"email warning: Your IP has changed\r\n")

err := smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("GMAIL"), to, msg)

if err != nil {

log.Fatal(err)

}

}