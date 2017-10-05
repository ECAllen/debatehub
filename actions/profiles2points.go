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
// Model: Singular (Profiles2point)
// DB Table: Plural (profiles2points)
// Resource: Plural (Profiles2points)
// Path: Plural (/profiles2points)
// View Template Folder: Plural (/templates/profiles2points/)

// Profiles2pointsResource is the resource for the profiles2point model
type Profiles2pointsResource struct {
	buffalo.Resource
}

// List gets all Profiles2points. This function is mapped to the path
// GET /profiles2points
func (v Profiles2pointsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	profiles2points := &models.Profiles2points{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(profiles2points)
	// to:
	// err := q.Order("created_at desc").All(profiles2points)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Profiles2points available inside the html template
	c.Set("profiles2points", profiles2points)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("profiles2points/index.html"))
}

// Show gets the data for one Profiles2point. This function is mapped to
// the path GET /profiles2points/{profiles2point_id}
func (v Profiles2pointsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Profiles2point
	profiles2point := &models.Profiles2point{}
	// To find the Profiles2point the parameter profiles2point_id is used.
	err := tx.Find(profiles2point, c.Param("profiles2point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make profiles2point available inside the html template
	c.Set("profiles2point", profiles2point)
	return c.Render(200, r.HTML("profiles2points/show.html"))
}

// New renders the form for creating a new Profiles2point.
// This function is mapped to the path GET /profiles2points/new
func (v Profiles2pointsResource) New(c buffalo.Context) error {
	// Make profiles2point available inside the html template
	c.Set("profiles2point", &models.Profiles2point{})
	return c.Render(200, r.HTML("profiles2points/new.html"))
}

// Create adds a Profiles2point to the DB. This function is mapped to the
// path POST /profiles2points
func (v Profiles2pointsResource) Create(c buffalo.Context) error {
	// Allocate an empty Profiles2point
	profiles2point := &models.Profiles2point{}
	// Bind profiles2point to the html form elements
	err := c.Bind(profiles2point)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(profiles2point)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make profiles2point available inside the html template
		c.Set("profiles2point", profiles2point)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("profiles2points/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Profiles2point was created successfully")
	// and redirect to the profiles2points index page
	return c.Redirect(302, "/profiles2points/%s", profiles2point.ID)
}

// Edit renders a edit form for a profiles2point. This function is
// mapped to the path GET /profiles2points/{profiles2point_id}/edit
func (v Profiles2pointsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Profiles2point
	profiles2point := &models.Profiles2point{}
	err := tx.Find(profiles2point, c.Param("profiles2point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make profiles2point available inside the html template
	c.Set("profiles2point", profiles2point)
	return c.Render(200, r.HTML("profiles2points/edit.html"))
}

// Update changes a profiles2point in the DB. This function is mapped to
// the path PUT /profiles2points/{profiles2point_id}
func (v Profiles2pointsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Profiles2point
	profiles2point := &models.Profiles2point{}
	err := tx.Find(profiles2point, c.Param("profiles2point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Profiles2point to the html form elements
	err = c.Bind(profiles2point)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(profiles2point)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make profiles2point available inside the html template
		c.Set("profiles2point", profiles2point)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("profiles2points/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Profiles2point was updated successfully")
	// and redirect to the profiles2points index page
	return c.Redirect(302, "/profiles2points/%s", profiles2point.ID)
}

// Destroy deletes a profiles2point from the DB. This function is mapped
// to the path DELETE /profiles2points/{profiles2point_id}
func (v Profiles2pointsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Profiles2point
	profiles2point := &models.Profiles2point{}
	// To find the Profiles2point the parameter profiles2point_id is used.
	err := tx.Find(profiles2point, c.Param("profiles2point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(profiles2point)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Profiles2point was destroyed successfully")
	// Redirect to the profiles2points index page
	return c.Redirect(302, "/profiles2points")
}
