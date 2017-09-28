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
// Model: Singular (Test)
// DB Table: Plural (tests)
// Resource: Plural (Tests)
// Path: Plural (/tests)
// View Template Folder: Plural (/templates/tests/)

// TestsResource is the resource for the test model
type TestsResource struct {
	buffalo.Resource
}

// List gets all Tests. This function is mapped to the path
// GET /tests
func (v TestsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	tests := &models.Tests{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(tests)
	// to:
	// err := q.Order("created_at desc").All(tests)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Tests available inside the html template
	c.Set("tests", tests)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("tests/index.html"))
}

// Show gets the data for one Test. This function is mapped to
// the path GET /tests/{test_id}
func (v TestsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Test
	test := &models.Test{}
	// To find the Test the parameter test_id is used.
	err := tx.Find(test, c.Param("test_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make test available inside the html template
	c.Set("test", test)
	return c.Render(200, r.HTML("tests/show.html"))
}

// New renders the formular for creating a new Test.
// This function is mapped to the path GET /tests/new
func (v TestsResource) New(c buffalo.Context) error {
	// Make test available inside the html template
	c.Set("test", &models.Test{})
	return c.Render(200, r.HTML("tests/new.html"))
}

// Create adds a Test to the DB. This function is mapped to the
// path POST /tests
func (v TestsResource) Create(c buffalo.Context) error {
	// Allocate an empty Test
	test := &models.Test{}
	// Bind test to the html form elements
	err := c.Bind(test)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(test)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make test available inside the html template
		c.Set("test", test)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("tests/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Test was created successfully")
	// and redirect to the tests index page
	return c.Redirect(302, "/tests/%s", test.ID)
}

// Edit renders a edit formular for a test. This function is
// mapped to the path GET /tests/{test_id}/edit
func (v TestsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Test
	test := &models.Test{}
	err := tx.Find(test, c.Param("test_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make test available inside the html template
	c.Set("test", test)
	return c.Render(200, r.HTML("tests/edit.html"))
}

// Update changes a test in the DB. This function is mapped to
// the path PUT /tests/{test_id}
func (v TestsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Test
	test := &models.Test{}
	err := tx.Find(test, c.Param("test_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Test to the html form elements
	err = c.Bind(test)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(test)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make test available inside the html template
		c.Set("test", test)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("tests/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Test was updated successfully")
	// and redirect to the tests index page
	return c.Redirect(302, "/tests/%s", test.ID)
}

// Destroy deletes a test from the DB. This function is mapped
// to the path DELETE /tests/{test_id}
func (v TestsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Test
	test := &models.Test{}
	// To find the Test the parameter test_id is used.
	err := tx.Find(test, c.Param("test_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(test)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Test was destroyed successfully")
	// Redirect to the tests index page
	return c.Redirect(302, "/tests")
}