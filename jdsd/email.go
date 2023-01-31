package jdsd

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
)

type EmailInfo struct {
	// the host of the email server,such as smtp.exmail.qq.com
	ServerHost string
	// the port of the email server,such as 465 in tengxun
	ServerPort int
	// the email address of sender
	FromEmail string
	// the password of the sender
	FromPassword string
	// the email address of receiver
	Recipient []string
}

// emailMessage defines the message of the email
var emailMessage *gomail.Message

func SendEmail(subject, body string, emailInfo *EmailInfo) error {
	if len(emailInfo.Recipient) == 0 {
		fmt.Println("收件人为空")
		return errors.New("收件人为空")
	}
	// init the email info
	emailMessage = gomail.NewMessage()
	// add the receiver for email
	emailMessage.SetHeader("To", emailInfo.Recipient...)
	// set the other name for the sender
	emailMessage.SetAddressHeader("From", emailInfo.FromEmail, "发邮件工具人")
	// set the subject
	emailMessage.SetHeader("Subject", subject)
	// set the content for the email
	emailMessage.SetBody("text/html", body)

	// create a new dialer
	d := gomail.NewDialer(emailInfo.ServerHost, emailInfo.ServerPort, emailInfo.FromEmail, emailInfo.FromPassword)
	err := d.DialAndSend(emailMessage)
	if err != nil {
		fmt.Println("邮件发送失败")
		return err
	}
	return nil
}
