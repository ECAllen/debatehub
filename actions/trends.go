package actions

import (
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Trend)
// DB Table: Plural (Trends)
// Resource: Plural (Trends)
// Path: Plural (/trends)
// View Template Folder: Plural (/templates/trends/)

// TrendsResource is the resource for the trend model
type TrendsResource struct {
	buffalo.Resource
}

// List gets all Trends. This function is mapped to the path
// GET /trends
func (v TrendsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	trends := &models.Trends{}
	// You can order your list here. Just change
	err := tx.All(trends)
	// to:
	// err := tx.Order("create_at desc").All(trends)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make trends available inside the html template
	c.Set("trends", trends)
	return c.Render(200, r.HTML("trends/index.html"))
}

// Show gets the data for one Trend. This function is mapped to
// the path GET /trends/{trend_id}
func (v TrendsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Trend
	trend := &models.Trend{}
	// To find the Trend the parameter trend_id is used.
	err := tx.Find(trend, c.Param("trend_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make trend available inside the html template
	c.Set("trend", trend)
	return c.Render(200, r.HTML("trends/show.html"))
}

// New renders the formular for creating a new trend.
// This function is mapped to the path GET /trends/new
func (v TrendsResource) New(c buffalo.Context) error {
	// Make trend available inside the html template
	c.Set("trend", &models.Trend{})
	return c.Render(200, r.HTML("trends/new.html"))
}

// Create adds a trend to the DB. This function is mapped to the
// path POST /trends
func (v TrendsResource) Create(c buffalo.Context) error {
	// Allocate an empty Trend
	trend := &models.Trend{}
	f := nulls.NewBool(false)
	trend.Publish = f
	trend.Reject = f

	// Bind trend to the html form elements
	err := c.Bind(trend)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(trend)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make trend available inside the html template
		c.Set("trend", trend)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("trends/submit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Trend was created successfully")
	// and redirect to the trends index page
	return c.Redirect(302, "/")
}

// Edit renders a edit formular for a trend. This function is
// mapped to the path GET /trends/{trend_id}/edit
func (v TrendsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Trend
	trend := &models.Trend{}
	err := tx.Find(trend, c.Param("trend_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make trend available inside the html template
	c.Set("trend", trend)
	return c.Render(200, r.HTML("trends/edit.html"))
}

// Update changes a trend in the DB. This function is mapped to
// the path PUT /trends/{trend_id}
func (v TrendsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Trend
	trend := &models.Trend{}
	err := tx.Find(trend, c.Param("trend_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind trend to the html form elements
	err = c.Bind(trend)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(trend)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make trend available inside the html template
		c.Set("trend", trend)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("trends/submit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Trend was updated successfully")
	// and redirect to the trends index page
	return c.Redirect(302, "/trends/%s", trend.ID)
}

// Destroy deletes a trend from the DB. This function is mapped
// to the path DELETE /trends/{trend_id}
func (v TrendsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Trend
	trend := &models.Trend{}
	// To find the Trend the parameter trend_id is used.
	err := tx.Find(trend, c.Param("trend_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(trend)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Trend was destroyed successfully")
	// Redirect to the trends index page
	return c.Redirect(302, "/trends")
}

// <================>Added<=================>

// New renders the formular for creating a new trend.
// This function is mapped to the path GET /trends/submit
func TrendsSubmit(c buffalo.Context) error {
	// Make trend available inside the html template
	c.Set("trend", &models.Trend{})
	// return c.Render(200, r.HTML("trends/new.html"))
	return c.Render(200, r.HTML("trends/submit.html"))
}

// List gets all Trends. This function is mapped to the path
// GET /trends
func TrendsAdmin(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	trends := &models.Trends{}
	// You can order your list here. Just change
	// err := tx.All(trends)
	// to:
	// err := tx.Order("create_at desc").All(trends)
	err := tx.Where("reject = false").Where("publish = false").All(trends)
	if err != nil {
		return errors.WithStack(err)
	}

	// Make trends available inside the html template
	c.Set("trends", trends)
	return c.Render(200, r.HTML("trends/admin.html"))
}