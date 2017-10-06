package actions

import (
	"github.com/ECAllen/debatehub/models"
	// "github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Ftree for forums
type Ftree struct {
	models.Thread
	models.Profile
	DebateID uuid.UUID
	Token    string
	Children []Ftree
}

func Thread(id uuid.UUID, tx *pop.Connection) (models.Thread, error) {

	// point used to hold the point
	thread := &models.Thread{}

	// Create query.
	q := tx.Where("ID = ?", id)

	// verify that the thread exists in
	// the Thread table
	exists, err := q.Exists(thread)
	if err != nil {
		return *thread, err
	}

	// Collect thread.
	if exists {
		err = q.First(thread)
		if err != nil {
			return *thread, err
		}
	}
	return *thread, err
}

func buildFTree(id uuid.UUID, debateID uuid.UUID, tx *pop.Connection, ftree *Ftree) error {

	// Get thread.
	thread, err := Thread(id, tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// Check if there is a profile associated
	// with this threada ID.
	profile2thread := &models.Profile2thread{}
	q := tx.Where("thread = ?", thread.ID)
	exists, err := q.Exists(profile2thread)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there is a profile then get it
	profile := models.Profile{}

	if exists {
		err = q.First(profile2thread)
		if err != nil {
			return errors.WithStack(err)
		}

		err = tx.Find(&profile, profile2thread.Profile)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	ftree.Thread = thread
	ftree.Profile = profile
	ftree.DebateID = debateID

	// check if this thread has any counterthreads
	thread2counterthread := &models.Thread2counterthread{}
	q = tx.Where("thread = ?", id)
	exists, err = q.Exists(thread2counterthread)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		threads2counterthreads := []models.Thread2counterthread{}
		err = q.All(&threads2counterthreads)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, c := range threads2counterthreads {
			var ft Ftree
			ft.Token = ftree.Token
			err = buildFTree(c.Counterthread, debateID, tx, &ft)
			if err != nil {
				return errors.WithStack(err)
			}
			ftree.Children = append(ftree.Children, ft)
		}
	}
	return errors.WithStack(err)
}

func buildFTreeRoot(id uuid.UUID, tx *pop.Connection, ftree *Ftree) error {

	// Assume there are child nodes or
	// would not have been called.
	q := tx.Where("debate = ?", id)

	// The debateThreads used to iterate through
	// to collect all the debate thread uuid's.
	debateThreads := []models.Debate2thread{}

	// Collect all the thread id's associated
	// with the debate.
	err := q.All(&debateThreads)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, dt := range debateThreads {
		var ft Ftree
		ft.Token = ftree.Token
		err = buildFTree(dt.Thread, id, tx, &ft)
		if err != nil {
			return errors.WithStack(err)
		}
		ftree.Children = append(ftree.Children, ft)
	}
	return nil
}
