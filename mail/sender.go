package mail

import (
	"bytes"
	"ditto/booking/config"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

//Mail -
type Mail struct {
	Auth    smtp.Auth
	Address string
	From    mail.Address
}

//New -
func New() *Mail {
	conf := config.Load()

	m := &Mail{
		Auth:    smtp.PlainAuth("", conf.Mail.Account, conf.Mail.Password, conf.Mail.Domain),
		Address: fmt.Sprintf("%v:%v", conf.Mail.Domain, conf.Mail.Port),
		From: mail.Address{
			Name:    conf.Mail.From.Name,
			Address: conf.Mail.From.Address,
		},
	}
	return m
}

//encodeSubject -
func (m *Mail) encodeSubject(subject string) string {
	var b bytes.Buffer
	strs := []string{}
	length := 13
	for k, c := range strings.Split(subject, "") {
		b.WriteString(c)
		if k%length == length-1 {
			strs = append(strs, b.String())
			b.Reset()
		}
	}
	if b.Len() > 0 {
		strs = append(strs, b.String())
	}
	// MIME エンコードする
	var bs bytes.Buffer
	bs.WriteString("Subject:")
	for _, line := range strs {
		bs.WriteString(" =?utf-8?B?")
		bs.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		bs.WriteString("?=\r\n")
	}
	return bs.String()
}

// 本文を 76 バイト毎に CRLF を挿入して返す
func (m *Mail) encodeBody(body string) string {
	b := bytes.NewBufferString(body)
	s := base64.StdEncoding.EncodeToString(b.Bytes())
	b2 := bytes.NewBuffer([]byte(""))
	for k, c := range strings.Split(s, "") {
		b2.WriteString(c)
		if k%76 == 75 {
			b2.WriteString("\r\n")
		}
	}
	return b2.String()
}

//Send -
func (m *Mail) Send(toaddr string, toname string, subject string, text string) error {
	to := mail.Address{
		Name:    toname,
		Address: toaddr,
	}
	var header bytes.Buffer
	header.WriteString("From: " + m.From.String() + "\r\n")
	header.WriteString("To: " + to.String() + "\r\n")
	header.WriteString(m.encodeSubject(subject))
	header.WriteString("MIME-Version: 1.0\r\n")
	header.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	header.WriteString("Content-Transfer-Encoding: base64\r\n")

	var body bytes.Buffer
	body.WriteString(text)

	var message bytes.Buffer
	message = header
	message.WriteString("\r\n")
	message.WriteString(base64.StdEncoding.EncodeToString(body.Bytes()))

	//send
	err := smtp.SendMail(m.Address, m.Auth, m.From.Address, []string{to.Address}, message.Bytes())
	if err != nil {
		return err
	}

	return nil
}
