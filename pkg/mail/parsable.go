package mail

import (
	"bytes"
	"errors"
	ht "html/template"
	"io"
	"strings"
	tt "text/template"
)

// Parsable represents an email message.
type Parsable struct {
	To          []string
	Subject     string
	CopyTo      []string         // the recipient is explicitly know about copy.
	BlindCopyTo []string         // the recipient is explicitly don't know about copy.
	Sender      string           // set another email sender.
	ReplyTo     string           // to whom the recipient will respond.
	Parts       []Part           // message body parts.
	PartValues  map[string]any   // used only with part body.
	Files       []*File          // message files.
	Settings    *ServiceSettings // advanced settings of the mailer service.
}

type ServiceSettings struct {
	Name   string // templateName, which will be read from db.
	Locale string // "en" or "ru". "en" by default.
}

func (p *Parsable) ToEmail(dsc *Email) *Email {
	email := dsc.SetSubject(p.Subject)

	// SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

	if p.Sender != "" {
		email.SetSender(p.Sender)
	}

	if p.ReplyTo != "" {
		email.SetReplyTo(p.ReplyTo)
	}

	if len(p.To) != 0 {
		email.AddTo(p.To...)
	}

	if len(p.CopyTo) != 0 {
		email.AddCc(p.CopyTo...)
	}

	if len(p.BlindCopyTo) != 0 {
		email.AddBcc(p.BlindCopyTo...)
	}

	for _, file := range p.Files {
		email.Attach(file)
	}

	email.Parts = p.Parts

	// insert template values
	for i, part := range email.Parts {
		var (
			t interface {
				Execute(wr io.Writer, data any) error
			}
			err error
		)

		switch part.ContentType {
		case TextHTML, TextAMP:
			t, err = ht.New("").Parse(string(part.Body))
		case TextPlain, TextCalendar:
			t, err = tt.New("").Parse(string(part.Body))
		default:
			email.Error = errors.New("content type is not found")
			return nil
		}
		if err != nil {
			email.Error = err
			return nil
		}

		buf := bytes.NewBuffer(make([]byte, 0, len(part.Body)))
		if err = t.Execute(buf, p.PartValues); err != nil {
			email.Error = err
			return nil
		}
		email.Parts[i].Body = buf.Bytes()
	}
	return email
}

func (p *Parsable) Recipients(delimiter string) string {
	sb := new(strings.Builder)
	for _, to := range p.To {
		sb.WriteString(to + delimiter)
	}
	for _, to := range p.CopyTo {
		sb.WriteString(to + delimiter)
	}
	for _, to := range p.BlindCopyTo {
		sb.WriteString(to + delimiter)
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()[:sb.Len()-len(delimiter)]
}
