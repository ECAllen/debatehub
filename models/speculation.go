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

type Speculation struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Summary     string     `json:"summary" db:"summary"`
	Description string     `json:"description" db:"description"`
	Publish     nulls.Bool `json:"publish" db:"publish"`
	Reject      nulls.Bool `json:"reject" db:"reject"`
}

// String is not required by pop and may be deleted
func (s Speculation) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Speculations is not required by pop and may be deleted
type Speculations []Speculation

// String is not required by pop and may be deleted
func (s Speculations) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (s *Speculation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Summary, Name: "Summary"},
		&validators.StringIsPresent{Field: s.Description, Name: "Description"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (s *Speculation) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (s *Speculation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
