package actions

import (
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"

	"fmt"
	"github.com/satori/go.uuid"
)

// DebatePagesResource is the resource for the debate model
type DebatePagesResource struct {
	buffalo.Resource
}

// List gets all Debates. This function is mapped to the path
// GET /debate_pages
func (v DebatePagesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	debates := &models.Debates{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(debates)
	// to:
	// err := q.Order("created_at desc").All(debate_pages)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Debates available inside the html template
	c.Set("debates", debates)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("debate_pages/index.html"))
}

// Show gets the data for one Debate. This function is mapped to
// the path GET /debate_pages/{debate_page_id}
func (v DebatePagesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	debate_id := c.Param("debate_page_id")

	tx := c.Value("tx").(*pop.Connection)

	// create the models
	debate := &models.Debate{}
	point := &models.Point{}
	debates2point := &models.Debates2point{}

	err := tx.Find(debate, debate_id)
	if err != nil {
		return errors.WithStack(err)
	}

	// check for existence counter points for debate
	q := tx.Where("debate = ?", debate_id)
	exists, err := q.Exists(debates2point)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		// collect counter points
		debatePoints := []models.Debates2point{}
		points := []models.Point{}

		// collect all the point id's for the debate
		err = q.All(&debatePoints)
		if err != nil {
			return errors.WithStack(err)
		}
		// iterate through the debate points
		for _, dp := range debatePoints {
			fmt.Println(dp.Point)

			// verify that the point exists in
			// the Point table
			qp := tx.Where("ID = ?", dp.Point)
			exists, err = qp.Exists(point)
			if err != nil {
				return errors.WithStack(err)
			}

			if exists {
				err = qp.All(&points)
				if err != nil {
					return errors.WithStack(err)

				}
				for _, pt := range points {
					fmt.Println(pt.Topic)
				}
			}
		}
		c.Set("points", points)
	}

	c.Set("point", point)
	c.Set("debate", debate)

	return c.Render(200, r.HTML("debate_pages/show.html"))
}

// New renders the formular for creating a new Debate.
// This function is mapped to the path GET /debate_pages/new
func (v DebatePagesResource) New(c buffalo.Context) error {
	// Make debate available inside the html template
	c.Set("debate", &models.Debate{})
	return c.Render(200, r.HTML("debate_pages/new.html"))
}

// Create adds a Debate to the DB. This function is mapped to the
// path POST /debate_pages
func (v DebatePagesResource) Create(c buffalo.Context) error {
	// Allocate an empty Debate
	debate := &models.Debate{}
	// Bind debate to the html form elements
	err := c.Bind(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make debate available inside the html template
		c.Set("debate", debate)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("debate_pages/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Debate was created successfully")
	// and redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages/%s", debate.ID)
}

// Edit renders a edit formular for a debate. This function is
// mapped to the path GET /debate_pages/{debate_page_id}/edit
func (v DebatePagesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make debate available inside the html template
	c.Set("debate", debate)
	return c.Render(200, r.HTML("debate_pages/edit.html"))
}

// Update changes a debate in the DB. This function is mapped to
// the path PUT /debate_pages/{debate_page_id}
func (v DebatePagesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Debate to the html form elements
	err = c.Bind(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make debate available inside the html template
		c.Set("debate", debate)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("debate_pages/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Debate was updated successfully")
	// and redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages/%s", debate.ID)
}

// Destroy deletes a debate from the DB. This function is mapped
// to the path DELETE /debate_pages/{debate_page_id}
func (v DebatePagesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	// To find the Debate the parameter debate_page_id is used.
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Debate was destroyed successfully")
	// Redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages")
}

func AddPoint(c buffalo.Context) error {
	debate_id := c.Param("debate_page_id")
	// Allocate an empty Point
	point := &models.Point{}
	point.Rank = 1
	// Bind point to the html form elements
	err := c.Bind(point)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(point)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make point available inside the html template
		c.Set("point", point)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		// Adding the user info to the session
		fmt.Println(verrs)
		c.Flash().Add("warning", fmt.Sprintf("%s", verrs))
		return c.Redirect(302, "/debate_pages/%s", debate_id)
	}

	// put point and debate into debate2points table
	debate2point := &models.Debates2point{}
	debate2point.Debate, err = uuid.FromString(debate_id)
	if err != nil {
		return errors.WithStack(err)
	}

	debate2point.Point = point.ID
	verrs, err = tx.ValidateAndCreate(debate2point)
	if err != nil {
		return errors.WithStack(err)
	}

	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_id)
}
