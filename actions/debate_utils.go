package actions

import (
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Points tree for debates
type Ptree struct {
	models.Point
	models.Profile
	DebateID uuid.UUID
	Token    string
	Children []Ptree
}

func Counterpoint2Point(id uuid.UUID, tx *pop.Connection) (models.Point, error) {

	parentPoint := models.Point{}

	p2c := &models.Points2counterpoint{}
	q := tx.Where("counterpoint = ?", id)
	exists, err := q.Exists(p2c)
	if err != nil {
		return parentPoint, errors.WithStack(err)
	}

	if exists {
		err = q.First(p2c)
		if err != nil {
			return parentPoint, errors.WithStack(err)
		}

		parentPoint, err = Point(p2c.Point, tx)
		if err != nil {
			return parentPoint, errors.WithStack(err)
		}
	}

	return parentPoint, err
}

func Point2Debate(id uuid.UUID, tx *pop.Connection) (models.Debate, error) {

	debate := models.Debate{}

	// If UUID os 0 then null model and just
	// return to avoind infinite loop.
	blank, err := uuid.FromString("00000000-0000-0000-0000-000000000000")
	if err != nil {
		return debate, errors.WithStack(err)
	}

	if uuid.Equal(id, blank) {
		return debate, nil
	}

	// Check if point is in debates to point table
	debates2point := &models.Debates2point{}
	q := tx.Where("point = ?", id)
	exists, err := q.Exists(debates2point)
	if err != nil {
		return debate, errors.WithStack(err)
	}

	if exists {
		err = q.First(debates2point)
		if err != nil {
			return debate, errors.WithStack(err)
		}

		// check for existence of debate in table
		debate, err = Debate(debates2point.Debate, tx)
		if err != nil {
			return debate, errors.WithStack(err)
		}
	} else {
		// Point is counter point.
		point, err := Counterpoint2Point(id, tx)
		if err != nil {
			return debate, errors.WithStack(err)
		}

		debate, err = Point2Debate(point.ID, tx)
		if err != nil {
			return debate, errors.WithStack(err)
		}
	}
	return debate, nil
}

func Debate(id uuid.UUID, tx *pop.Connection) (models.Debate, error) {

	// point used to hold the point
	debate := &models.Debate{}

	// Create query.
	q := tx.Where("ID = ?", id)

	// verify that the debate exists in
	// the Debate table
	exists, err := q.Exists(debate)
	if err != nil {
		return *debate, err
	}

	// Collect debate.
	if exists {
		err = q.First(debate)
		if err != nil {
			return *debate, err
		}
	}
	return *debate, err
}

func Point(id uuid.UUID, tx *pop.Connection) (models.Point, error) {

	// point used to hold the point
	point := &models.Point{}

	// Create query.
	q := tx.Where("ID = ?", id)

	// verify that the point exists in
	// the Point table
	exists, err := q.Exists(point)
	if err != nil {
		return *point, err
	}

	// Collect point.
	if exists {
		err = q.First(point)
		if err != nil {
			return *point, err
		}
	}
	return *point, err
}

func buildTree(id uuid.UUID, debateID uuid.UUID, tx *pop.Connection, ptree *Ptree) error {

	// Get point.
	point, err := Point(id, tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// Check if there is a profile associated
	// with this point ID.
	profile2point := &models.Profiles2point{}
	q := tx.Where("point = ?", point.ID)
	exists, err := q.Exists(profile2point)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there is a profile then get it
	profile := models.Profile{}

	if exists {
		err = q.First(profile2point)
		if err != nil {
			return errors.WithStack(err)
		}

		err = tx.Find(&profile, profile2point.Profile)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	ptree.Point = point
	ptree.Profile = profile
	ptree.DebateID = debateID

	// check if this point has any counterpoints
	p2c := &models.Points2counterpoint{}
	q = tx.Where("point = ?", id)
	exists, err = q.Exists(p2c)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		p2cs := []models.Points2counterpoint{}
		err = q.All(&p2cs)
		if err != nil {
			return errors.WithStack(err)
		}

		// Iterate through each counter point
		// create a ptree node recusrsivley call
		// buildtree on node.
		for _, c := range p2cs {
			var pt Ptree
			pt.Token = ptree.Token
			err = buildTree(c.Counterpoint, debateID, tx, &pt)
			if err != nil {
				return errors.WithStack(err)
			}
			ptree.Children = append(ptree.Children, pt)
		}

	}
	return errors.WithStack(err)
}

func buildTreeRoot(id uuid.UUID, tx *pop.Connection, ptree *Ptree) error {

	// Check if this debate has any points.
	debatePoint := &models.Debates2point{}
	q := tx.Where("debate = ?", id)
	exists, err := q.Exists(debatePoint)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		// The debatePoints used to iterate through
		// to collect all the debate point uuid's.
		debatePoints := []models.Debates2point{}

		// Collect all the point id's associated
		// with the debate.
		err = q.All(&debatePoints)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, dp := range debatePoints {
			var pt Ptree
			pt.Token = ptree.Token
			err = buildTree(dp.Point, id, tx, &pt)
			if err != nil {
				return errors.WithStack(err)
			}
			ptree.Children = append(ptree.Children, pt)
		}
	}
	return errors.WithStack(err)
}

func insertProfile2Point(pointID uuid.UUID, tx *pop.Connection, c buffalo.Context) error {

	// Associate profile with debate
	profile2point := &models.Profiles2point{}

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err := errors.New("should not have gotten here, check authentication path")
		return errors.WithStack(err)
	} else {
		profile2point.Profile = userID.(uuid.UUID)
	}

	profile2point.Point = pointID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(profile2point)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}
	return nil
}
