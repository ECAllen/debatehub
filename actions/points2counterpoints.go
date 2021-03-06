package actions

import (
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Points2counterpoint)
// DB Table: Plural (points2counterpoints)
// Resource: Plural (Points2counterpoints)
// Path: Plural (/points2counterpoints)
// View Template Folder: Plural (/templates/points2counterpoints/)

// Points2counterpointsResource is the resource for the points2counterpoint model
type Points2counterpointsResource struct {
	buffalo.Resource
}

// List gets all Points2counterpoints. This function is mapped to the path
// GET /points2counterpoints
func (v Points2counterpointsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	points2counterpoints := &models.Points2counterpoints{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(points2counterpoints)
	// to:
	// err := q.Order("created_at desc").All(points2counterpoints)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Points2counterpoints available inside the html template
	c.Set("points2counterpoints", points2counterpoints)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("points2counterpoints/index.html"))
}

// Show gets the data for one Points2counterpoint. This function is mapped to
// the path GET /points2counterpoints/{points2counterpoint_id}
func (v Points2counterpointsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Points2counterpoint
	points2counterpoint := &models.Points2counterpoint{}
	// To find the Points2counterpoint the parameter points2counterpoint_id is used.
	err := tx.Find(points2counterpoint, c.Param("points2counterpoint_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make points2counterpoint available inside the html template
	c.Set("points2counterpoint", points2counterpoint)
	return c.Render(200, r.HTML("points2counterpoints/show.html"))
}

// New renders the formular for creating a new Points2counterpoint.
// This function is mapped to the path GET /points2counterpoints/new
func (v Points2counterpointsResource) New(c buffalo.Context) error {
	// Make points2counterpoint available inside the html template
	c.Set("points2counterpoint", &models.Points2counterpoint{})
	return c.Render(200, r.HTML("points2counterpoints/new.html"))
}

// Create adds a Points2counterpoint to the DB. This function is mapped to the
// path POST /points2counterpoints
func (v Points2counterpointsResource) Create(c buffalo.Context) error {
	// Allocate an empty Points2counterpoint
	points2counterpoint := &models.Points2counterpoint{}
	// Bind points2counterpoint to the html form elements
	err := c.Bind(points2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(points2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make points2counterpoint available inside the html template
		c.Set("points2counterpoint", points2counterpoint)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("points2counterpoints/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Points2counterpoint was created successfully")
	// and redirect to the points2counterpoints index page
	return c.Redirect(302, "/points2counterpoints/%s", points2counterpoint.ID)
}

// Edit renders a edit formular for a points2counterpoint. This function is
// mapped to the path GET /points2counterpoints/{points2counterpoint_id}/edit
func (v Points2counterpointsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Points2counterpoint
	points2counterpoint := &models.Points2counterpoint{}
	err := tx.Find(points2counterpoint, c.Param("points2counterpoint_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make points2counterpoint available inside the html template
	c.Set("points2counterpoint", points2counterpoint)
	return c.Render(200, r.HTML("points2counterpoints/edit.html"))
}

// Update changes a points2counterpoint in the DB. This function is mapped to
// the path PUT /points2counterpoints/{points2counterpoint_id}
func (v Points2counterpointsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Points2counterpoint
	points2counterpoint := &models.Points2counterpoint{}
	err := tx.Find(points2counterpoint, c.Param("points2counterpoint_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Points2counterpoint to the html form elements
	err = c.Bind(points2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(points2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make points2counterpoint available inside the html template
		c.Set("points2counterpoint", points2counterpoint)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("points2counterpoints/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Points2counterpoint was updated successfully")
	// and redirect to the points2counterpoints index page
	return c.Redirect(302, "/points2counterpoints/%s", points2counterpoint.ID)
}

// Destroy deletes a points2counterpoint from the DB. This function is mapped
// to the path DELETE /points2counterpoints/{points2counterpoint_id}
func (v Points2counterpointsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Points2counterpoint
	points2counterpoint := &models.Points2counterpoint{}
	// To find the Points2counterpoint the parameter points2counterpoint_id is used.
	err := tx.Find(points2counterpoint, c.Param("points2counterpoint_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(points2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Points2counterpoint was destroyed successfully")
	// Redirect to the points2counterpoints index page
	return c.Redirect(302, "/points2counterpoints")
}
