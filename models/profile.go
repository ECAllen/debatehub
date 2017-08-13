package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Profile struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	FirstName   string       `json:"firstname" db:"firstname"`
	LastName    string       `json:"lastname" db:"lastname"`
	Provider    string       `json:"provider" db:"provider"`
	UserID      string       `json:"userid" db:"userid"`
	Email       string       `json:"email" db:"email"`
	NickName    string       `json:"nickname" db:"nickname"`
	Location    nulls.String `json:"location" db:"location"`
	AvatarURL   nulls.String `json:"avatarurl" db:"avatarurl"`
	Description nulls.String `json:"description" db:"description"`
}

// String is not required by pop and may be deleted
func (p Profile) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Profiles is not required by pop and may be deleted
type Profiles []Profile

// String is not required by pop and may be deleted
func (p Profiles) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (p *Profile) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: p.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: p.Email, Name: "Email"},
		&validators.StringIsPresent{Field: p.NickName, Name: "NickName"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (p *Profile) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (p *Profile) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
