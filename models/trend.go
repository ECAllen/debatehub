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

type Trend struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Summary     string       `json:"summary" db:"summary"`
	Description string       `json:"description" db:"description"`
	Thumbnail   nulls.String `json:"thumbnail" db:"thumbnail"`
	Publish     nulls.Bool   `json:"publish" db:"publish"`
	Reject      nulls.Bool   `json:"reject" db:"reject"`
}

// String is not required by pop and may be deleted
func (t Trend) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Trends is not required by pop and may be deleted
type Trends []Trend

// String is not required by pop and may be deleted
func (t Trends) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (t *Trend) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Summary, Name: "Summary"},
		&validators.StringIsPresent{Field: t.Description, Name: "Description"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (t *Trend) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (t *Trend) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
