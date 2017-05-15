package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Email struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Email     string    `json:"email" db:"email"`
}

// String is not required by pop and may be deleted
func (e Email) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Emails is not required by pop and may be deleted
type Emails []Email

// String is not required by pop and may be deleted
func (e Emails) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (e *Email) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: e.Email, Name: "Email"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (e *Email) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (e *Email) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
