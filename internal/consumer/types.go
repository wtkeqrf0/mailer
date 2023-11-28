package consumer

// Email represents an email message.
type Email struct {
	To          []string `bson:"-"`
	Subject     string
	CopyTo      []string        `json:",omitempty"`          // the recipient is explicitly know about copy.
	BlindCopyTo []string        `json:",omitempty"`          // the recipient is explicitly don't know about copy.
	Sender      string          `json:",omitempty"`          // set another email sender.
	ReplyTo     string          `json:",omitempty"`          // to whom the recipient will respond.
	Parts       []Part          `json:",omitempty"`          // message body parts.
	PartValues  map[string]any  `json:",omitempty" bson:"-"` // used only with part body.
	Files       []*File         `json:",omitempty"`          // message files.
	Settings    ServiceSettings `json:",omitempty" bson:"-"` // advanced settings of the mailer service.
}

// Part represents the different content parts of an email body.
//
// At first, tries to find template by Name, and fill it by Value field.
// If name is not found, given Body and ContentType should be used.
type Part struct {
	ContentType ContentType `json:",omitempty"`
	Body        string      `json:",omitempty"`
}

// File represents the file that can be added to the email message.
// You can add attachment from file in path, from base64 string or from []byte.
// You can define if attachment is inline or not.
// Only one, Data, B64Data or FilePath is supported. If multiple are set, then
// the first in that order is used.
type File struct {
	// FilePath is the path of the file to attach.
	FilePath string `json:",omitempty"`
	// Name is the name of file in attachment. Required for Data and B64Data. Optional for FilePath.
	Name string `json:",omitempty"`
	// MimeType of attachment. If empty then is obtained from Name (if not empty) or FilePath. If cannot obtained, application/octet-stream is set.
	MimeType string `json:",omitempty"`
	// B64Data is the base64 string to attach.
	B64Data string `json:"b64Data,omitempty"`
	// Data is the []byte of file to attach.
	Data []byte `json:",omitempty"`
	// Inline defines if attachment is inline or not.
	Inline bool `json:",omitempty"`
}

type ContentType string

const (
	// TextPlain sets body type to text/plain in message body.
	TextPlain ContentType = "text/plain"
	// TextHTML sets body type to text/html in message body.
	TextHTML ContentType = "text/html"
	// TextCalendar sets body type to text/calendar in message body.
	TextCalendar ContentType = "text/calendar"
	// TextAMP sets body type to text/x-amp-html in message body.
	TextAMP ContentType = "text/x-amp-html"
)

type ServiceSettings struct {
	Name   TemplateName `json:",omitempty"` // templateName, which will be read from db.
	From   Address      `json:",omitempty"` // address, from which email will be sent. SingleSender by default.
	Locale Locale       `json:",omitempty"` // "en" or "ru". "en" by default.
}

type Address string

const (
	SingleSenderAddress Address = "singleSender"
	MailingAddress      Address = "mailing"
)

type TemplateName string

// Part possible template names.
const (
	TemplateHello          TemplateName = "hello"
	TemplateConfirmEmail   TemplateName = "confirm_email"
	TemplateChangePassword TemplateName = "change_password"
)

// Locale can identify language for the template.
type Locale string

const (
	LocaleRu Locale = "ru"
	LocaleEn Locale = "en"
)
