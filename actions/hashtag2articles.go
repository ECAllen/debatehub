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
// Model: Singular (Hashtag2article)
// DB Table: Plural (hashtag2articles)
// Resource: Plural (Hashtag2articles)
// Path: Plural (/hashtag2articles)
// View Template Folder: Plural (/templates/hashtag2articles/)

// Hashtag2articlesResource is the resource for the hashtag2article model
type Hashtag2articlesResource struct {
	buffalo.Resource
}

// List gets all Hashtag2articles. This function is mapped to the path
// GET /hashtag2articles
func (v Hashtag2articlesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	hashtag2articles := &models.Hashtag2articles{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(hashtag2articles)
	// to:
	// err := q.Order("created_at desc").All(hashtag2articles)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Hashtag2articles available inside the html template
	c.Set("hashtag2articles", hashtag2articles)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("hashtag2articles/index.html"))
}

// Show gets the data for one Hashtag2article. This function is mapped to
// the path GET /hashtag2articles/{hashtag2article_id}
func (v Hashtag2articlesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Hashtag2article
	hashtag2article := &models.Hashtag2article{}
	// To find the Hashtag2article the parameter hashtag2article_id is used.
	err := tx.Find(hashtag2article, c.Param("hashtag2article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make hashtag2article available inside the html template
	c.Set("hashtag2article", hashtag2article)
	return c.Render(200, r.HTML("hashtag2articles/show.html"))
}

// New renders the form for creating a new Hashtag2article.
// This function is mapped to the path GET /hashtag2articles/new
func (v Hashtag2articlesResource) New(c buffalo.Context) error {
	// Make hashtag2article available inside the html template
	c.Set("hashtag2article", &models.Hashtag2article{})
	return c.Render(200, r.HTML("hashtag2articles/new.html"))
}

// Create adds a Hashtag2article to the DB. This function is mapped to the
// path POST /hashtag2articles
func (v Hashtag2articlesResource) Create(c buffalo.Context) error {
	// Allocate an empty Hashtag2article
	hashtag2article := &models.Hashtag2article{}
	// Bind hashtag2article to the html form elements
	err := c.Bind(hashtag2article)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(hashtag2article)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make hashtag2article available inside the html template
		c.Set("hashtag2article", hashtag2article)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("hashtag2articles/new.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Hashtag2article was created successfully")
	// and redirect to the hashtag2articles index page
	return c.Redirect(302, "/hashtag2articles/%s", hashtag2article.ID)
}

// Edit renders a edit form for a hashtag2article. This function is
// mapped to the path GET /hashtag2articles/{hashtag2article_id}/edit
func (v Hashtag2articlesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Hashtag2article
	hashtag2article := &models.Hashtag2article{}
	err := tx.Find(hashtag2article, c.Param("hashtag2article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make hashtag2article available inside the html template
	c.Set("hashtag2article", hashtag2article)
	return c.Render(200, r.HTML("hashtag2articles/edit.html"))
}

// Update changes a hashtag2article in the DB. This function is mapped to
// the path PUT /hashtag2articles/{hashtag2article_id}
func (v Hashtag2articlesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Hashtag2article
	hashtag2article := &models.Hashtag2article{}
	err := tx.Find(hashtag2article, c.Param("hashtag2article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Hashtag2article to the html form elements
	err = c.Bind(hashtag2article)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(hashtag2article)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make hashtag2article available inside the html template
		c.Set("hashtag2article", hashtag2article)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("hashtag2articles/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Hashtag2article was updated successfully")
	// and redirect to the hashtag2articles index page
	return c.Redirect(302, "/hashtag2articles/%s", hashtag2article.ID)
}

// Destroy deletes a hashtag2article from the DB. This function is mapped
// to the path DELETE /hashtag2articles/{hashtag2article_id}
func (v Hashtag2articlesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Hashtag2article
	hashtag2article := &models.Hashtag2article{}
	// To find the Hashtag2article the parameter hashtag2article_id is used.
	err := tx.Find(hashtag2article, c.Param("hashtag2article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(hashtag2article)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Hashtag2article was destroyed successfully")
	// Redirect to the hashtag2articles index page
	return c.Redirect(302, "/hashtag2articles")
}