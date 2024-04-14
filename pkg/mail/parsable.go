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
	Name   Template // templateName, which will be read from db.
	Locale string   // "en" or "ru". "en" by default.
}

type Template uint8

// possible template values.
const (
	TemplateHello Template = iota
	TemplateConfirmEmail
	TemplateChangePassword
)

func (pe *Parsable) ToEmail(dsc *Email) *Email {
	email := dsc.SetSubject(pe.Subject)

	// SetDSN([]mail.DSN{mail.SUCCESS, mail.FAILURE}, false)

	if pe.Sender != "" {
		email.SetSender(pe.Sender)
	}

	if pe.ReplyTo != "" {
		email.SetReplyTo(pe.ReplyTo)
	}

	if len(pe.To) != 0 {
		email.AddTo(pe.To...)
	}

	if len(pe.CopyTo) != 0 {
		email.AddCc(pe.CopyTo...)
	}

	if len(pe.BlindCopyTo) != 0 {
		email.AddBcc(pe.BlindCopyTo...)
	}

	for _, file := range pe.Files {
		email.Attach(file)
	}

	email.Parts = pe.Parts

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
		if err = t.Execute(buf, pe.PartValues); err != nil {
			email.Error = err
			return nil
		}
		email.Parts[i].Body = buf.Bytes()
	}
	return email
}

func (pe *Parsable) Recipients(delimiter string) string {
	sb := new(strings.Builder)
	for _, to := range pe.To {
		sb.WriteString(to + delimiter)
	}
	for _, to := range pe.CopyTo {
		sb.WriteString(to + delimiter)
	}
	for _, to := range pe.BlindCopyTo {
		sb.WriteString(to + delimiter)
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()[:sb.Len()-len(delimiter)]
}
