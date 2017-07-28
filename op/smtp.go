package op

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"bytes"

	"github.com:weisd/zerolog"
)

const (
	subjectPhrase = "Diagnostic message from server"
)

var _ zerolog.LevelWriter = &SmtpWriter{}

// smtpWriter implements LoggerInterface and is used to send emails via given SMTP-server.
type SmtpWriter struct {
	Username           string        `json:"Username"`
	Password           string        `json:"password"`
	Host               string        `json:"Host"`
	Subject            string        `json:"subject"`
	RecipientAddresses []string      `json:"sendTos"`
	Level              zerolog.Level `json:"level"`

	sender chan []byte
}

// create smtp writer.
func NewSmtpWriter(username, password, host, subject string, sendto []string, level zerolog.Level) zerolog.LevelWriter {
	s := &SmtpWriter{
		Username:           username,
		Password:           password,
		Host:               host,
		Subject:            subject,
		RecipientAddresses: sendto,
		Level:              level,
		sender:             make(chan []byte, 1000),
	}

	go s.loop()

	return s
}

func (s *SmtpWriter) loop() {
	for {
		select {
		case msg := <-s.sender:
			_, err := s.Write(msg)
			if err != nil {
				fmt.Println("smtp send err ", err)
			}
		}
	}
}

// init smtp writer with json config.
// config like:
//	{
//		"Username":"example@gmail.com",
//		"password:"password",
//		"host":"smtp.gmail.com:465",
//		"subject":"email title",
//		"sendTos":["email1","email2"],
//		"level":LevelError
//	}
// func (s *SmtpWriter) Init(jsonconfig string) error {
// 	return json.Unmarshal([]byte(jsonconfig), sw)
// }

func (s *SmtpWriter) Write(b []byte) (int, error) {
	hp := strings.Split(s.Host, ":")

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		s.Username,
		s.Password,
		hp[0],
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.

	var buf bytes.Buffer
	buf.WriteString("To: ")
	buf.WriteString(strings.Join(s.RecipientAddresses, ";"))
	buf.WriteString("\r\nFrom: ")
	buf.WriteString(s.Username)
	buf.WriteString(">\r\nSubject: ")
	buf.WriteString(s.Subject)
	buf.WriteString("\r\n")
	buf.WriteString("Content-Type: text/plain" + "; charset=UTF-8")
	buf.WriteString("\r\n\r\n")

	buf.WriteString(fmt.Sprintf("%s: ", time.Now().Format("2006-01-02 15:04:05")))
	buf.Write(b)

	// mailmsg := []byte("To: " + strings.Join(s.RecipientAddresses, ";") + "\r\nFrom: " + s.Username + "<" + s.Username +
	// 	">\r\nSubject: " + s.Subject + "\r\n" + content_type + "\r\n\r\n" + fmt.Sprintf(".%s", time.Now().Format("2006-01-02 15:04:05")) + msg)

	return len(b), smtp.SendMail(
		s.Host,
		auth,
		s.Username,
		s.RecipientAddresses,
		buf.Bytes(),
	)
}

// write message in smtp writer.
// it will send an email with subject and only this message.
func (s *SmtpWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l < s.Level {
		return len(p), nil
	}

	return s.Write(p)
	// go func() {
	// 	s.sender <- p
	// }()
	// return len(p), nil
}
