package actions

import (
	"fmt"
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
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

func insertProfile2Thread(threadID uuid.UUID, tx *pop.Connection, c buffalo.Context) error {

	// Associate profile with debate
	profile2thread := &models.Profile2thread{}

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err := errors.New("should not have gotten here, check authentication path")
		return errors.WithStack(err)
	} else {
		profile2thread.Profile = userID.(uuid.UUID)
	}

	profile2thread.Thread = threadID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(profile2thread)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}
	return nil
}

func AddThread(c buffalo.Context) error {
	// ==================================
	// Pull out params

	// The debate page id is needed in case we need to redirect
	// debate page if errors.
	debate_page_id := c.Param("debate_page_id")
	// Point_id is the existing point which this "counter point"
	// is attached.
	parent_thread_id := c.Param("parent_thread_id")

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// ==================================
	// Create the new thread.
	// ==================================
	newthread := &models.Thread{}
	newthread.Rank = 1

	// Bind thread to the html form elements
	err := c.Bind(newthread)
	if err != nil {
		return errors.WithStack(err)
	}

	// Validate the data from the html form.
	verrs, err := tx.ValidateAndCreate(newthread)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make point available inside the html template.
		c.Set("newthread", newthread)
		// Make the errors available inside the html template.
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		c.Flash().Add("warning", fmt.Sprintf("%s", verrs))
		// Redirect to the original debate page where the
		// counter point was created.
		return c.Redirect(422, "/debate_pages/%s", debate_page_id)
	}

	// If the debate ID and thread ID are the same
	// then add entry into debate2threads table
	// otherwise add entry into thread2counterthreads
	// table.
	if debate_page_id == parent_thread_id {
		debate2thread := &models.Debate2thread{}
		debate2thread.Debate, err = uuid.FromString(debate_page_id)
		if err != nil {
			return errors.WithStack(err)
		}

		// add point id
		debate2thread.Thread = newthread.ID
		verrs, err = tx.ValidateAndCreate(debate2thread)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		// Put new thread and parent thread into
		// thread2counterthreads table.
		thread2counterthread := &models.Thread2counterthread{}

		thread2counterthread.Thread, err = uuid.FromString(parent_thread_id)
		if err != nil {
			return errors.WithStack(err)
		}

		thread2counterthread.Counterthread = newthread.ID
		verrs, err = tx.ValidateAndCreate(thread2counterthread)
		if err != nil {
			return errors.WithStack(err)
		}

		insertProfile2Thread(newthread.ID, tx, c)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_page_id)
}
