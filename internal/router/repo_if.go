// Code generated by ifacemaker; DO NOT EDIT.

package router

import (
	"mailer/pkg/mail"
)

// Repository ...
type Repository interface {
	// GetTemplateByName from the db, and save it into given part, if it can be found by name.
	GetTemplateByName(email *mail.Parsable) error
}
