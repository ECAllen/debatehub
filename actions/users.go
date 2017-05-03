package actions

import (
	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (Users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

// UsersResource is the resource for the user model
type UsersResource struct {
	buffalo.Resource
}

// List gets all Users. This function is mapped to the the path
// GET /users
func (v UsersResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	users := &models.Users{}
	// You can order your list here. Just change
	err := tx.All(users)
	// to:
	// err := tx.Order("(case when completed then 1 else 2 end) desc, lower([sort_parameter]) asc").All(users)
	// Don't forget to change [sort_parameter] to the parameter of
	// your model, which should be used for sorting.
	if err != nil {
		return err
	}
	// Make users available inside the html template
	c.Set("users", users)
	return c.Render(200, r.HTML("users/index.html"))
}

// Show gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
func (v UsersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty User
	user := &models.User{}
	// To find the User the parameter user_id is used.
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	// Make user available inside the html template
	c.Set("user", user)
	return c.Render(200, r.HTML("users/show.html"))
}

// New renders the formular for creating a new user.
// This function is mapped to the path GET /users/new
func (v UsersResource) New(c buffalo.Context) error {
	// Make user available inside the html template
	c.Set("user", &models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

// Create adds a user to the DB. This function is mapped to the
// path POST /users
func (v UsersResource) Create(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}
	// Bind user to the html form elements
	err := c.Bind(user)
	if err != nil {
		return err
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make user available inside the html template
		c.Set("user", user)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("users/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "User was created successfully")
	// and redirect to the users index page
	return c.Redirect(302, "/users/%s", user.ID)
}

// Edit renders a edit formular for a user. This function is
// mapped to the path GET /users/{user_id}/edit
func (v UsersResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty User
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	// Make user available inside the html template
	c.Set("user", user)
	return c.Render(200, r.HTML("users/edit.html"))
}

// Update changes a user in the DB. This function is mapped to
// the path PUT /users/{user_id}
func (v UsersResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty User
	user := &models.User{}
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	// Bind user to the html form elements
	err = c.Bind(user)
	if err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		// Make user available inside the html template
		c.Set("user", user)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("users/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "User was updated successfully")
	// and redirect to the users index page
	return c.Redirect(302, "/users/%s", user.ID)
}

// Destroy deletes a user from the DB. This function is mapped
// to the path DELETE /users/{user_id}
func (v UsersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty User
	user := &models.User{}
	// To find the User the parameter user_id is used.
	err := tx.Find(user, c.Param("user_id"))
	if err != nil {
		return err
	}
	err = tx.Destroy(user)
	if err != nil {
		return err
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "User was destroyed successfully")
	// Redirect to the users index page
	return c.Redirect(302, "/users")
}
