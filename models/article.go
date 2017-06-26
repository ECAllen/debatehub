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

type Article struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Url         string       `json:"url" db:"url"`
	Description nulls.String `json:"description" db:"description"`
	Headline    nulls.String `json:"headline" db:"headline"`
	Type        string       `json:"type" db:"type"`
	Thumbnail   nulls.String `json:"thumbnail" db:"thumbnail"`
	Publish     nulls.Bool   `json:"publish" db:"publish"`
	Reject      nulls.Bool   `json:"reject" db:"reject"`
}

// String is not required by pop and may be deleted
func (a Article) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Articles is not required by pop and may be deleted
type Articles []Article

// String is not required by pop and may be deleted
func (a Articles) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (a *Article) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Url, Name: "Url"},
		&validators.StringIsPresent{Field: a.Type, Name: "Type"},
	), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (a *Article) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (a *Article) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
