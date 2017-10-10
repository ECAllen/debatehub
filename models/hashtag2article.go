package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

type Hashtag2article struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Hashtag   uuid.UUID `json:"hashtag" db:"hashtag"`
	Article   uuid.UUID `json:"article" db:"article"`
}

// String is not required by pop and may be deleted
func (h Hashtag2article) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Hashtag2articles is not required by pop and may be deleted
type Hashtag2articles []Hashtag2article

// String is not required by pop and may be deleted
func (h Hashtag2articles) String() string {
	jh, _ := json.Marshal(h)
	return string(jh)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (h *Hashtag2article) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (h *Hashtag2article) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (h *Hashtag2article) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
