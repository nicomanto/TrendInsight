package config

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

const (
	// error mail key
	MAIL_CONFIG_ROOT    = "/trend_insight/mail/"
	MAIL_SENDER_KEY     = "sender"
	MAIL_SENDER_PWD_KEY = "sender_pwd"
	MAIL_RECIPIENTS_KEY = "recipiens"
)

type MailClientConfigParam struct {
	Sender    string
	SenderPwd string
}

type GoMail struct {
	dialer *gomail.Dialer
	msg    *gomail.Message
}

func NewMailClientConfig(parameterStore ParameterStore) (*MailClientConfigParam, error) {
	mailClientConfig, err := parameterStore.GetParametersByPath(MAIL_CONFIG_ROOT, false, 3)
	if err != nil {
		return nil, err
	}
	config := MailClientConfigParam{}
	if val, ok := mailClientConfig[MAIL_SENDER_KEY]; ok {
		config.Sender = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", MAIL_SENDER_KEY)
	}
	if val, ok := mailClientConfig[MAIL_SENDER_PWD_KEY]; ok {
		config.SenderPwd = val
	} else {
		return nil, fmt.Errorf("failed to find parameter %s", MAIL_SENDER_PWD_KEY)
	}
	return &config, nil
}

// NewMailMsgAndDialer setup email dialer and sender of everey email messages
func NewSetupMailMsgAndDialer(cfgParam MailClientConfigParam) *GoMail {
	// setup default email message
	msg := gomail.NewMessage()
	// set E-Mail sender
	msg.SetHeader("From", cfgParam.Sender)

	// setup dialer
	// Settings for SMTP server
	dialer := gomail.NewDialer("smtp.gmail.com", 587, cfgParam.Sender, cfgParam.SenderPwd)
	// this is only needed when SSL/TLS certificate is not valid on server.
	// in production this should be set to false.
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &GoMail{
		dialer: dialer,
		msg:    msg,
	}
}

// SendErrorMail send email to given recipiens with the given error msg
func (gc *GoMail) SendErrorMail(recipiens []string, errorMsg string) {
	// add timestamp to text
	errorMsg = errorMsg + "\n\nDetected at " + time.Now().UTC().Format(time.RFC822)
	// set recipiens
	addresses := make([]string, len(recipiens))
	for i := range addresses {
		addresses[i] = gc.msg.FormatAddress(recipiens[i], "")
	}
	gc.msg.SetHeader("To", addresses...)
	// set E-Mail subject
	gc.msg.SetHeader("Subject", "Trend insight bot error")
	// et E-Mail body. You can set plain text or html with text/html
	gc.msg.SetBody("text/plain", errorMsg)
	// Now send E-Mail
	if err := gc.dialer.DialAndSend(gc.msg); err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("Error email sent")
	}
}
