package support

import (
	"crypto/tls"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var dialer *gomail.Dialer
var msg *gomail.Message

// SetupMailMsgAndDialer setup email dialer and sender of everey email messages
func SetupMailMsgAndDialer(sender string, senderPwd string) {
	// setup default email message
	msg = gomail.NewMessage()
	// set E-Mail sender
	msg.SetHeader("From", sender)

	// setup dialer
	// Settings for SMTP server
	dialer = gomail.NewDialer("smtp.gmail.com", 587, sender, senderPwd)
	// this is only needed when SSL/TLS certificate is not valid on server.
	// in production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}

// SendErrorMail send email to given recipiens with the given error msg
func SendErrorMail(recipiens []string, errorMsg string) {
	// add timestamp to text
	errorMsg = errorMsg + "\n\nDetected at " + time.Now().UTC().Format(time.RFC822)
	// set recipiens
	addresses := make([]string, len(recipiens))
	for i := range addresses {
		addresses[i] = msg.FormatAddress(recipiens[i], "")
	}
	msg.SetHeader("To", addresses...)
	// set E-Mail subject
	msg.SetHeader("Subject", "Trend insight bot error")
	// et E-Mail body. You can set plain text or html with text/html
	msg.SetBody("text/plain", errorMsg)
	// Now send E-Mail
	if err := dialer.DialAndSend(msg); err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("Error email sent")
	}
}
