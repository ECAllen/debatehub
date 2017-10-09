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
// Model: Singular (Article)
// DB Table: Plural (Articles)
// Resource: Plural (Articles)
// Path: Plural (/articles)
// View Template Folder: Plural (/templates/articles/)

// ArticlesResource is the resource for the article model
type ArticlesResource struct {
	buffalo.Resource
}

// List gets all Articles. This function is mapped to the path
// GET /articles
func (v ArticlesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	articles := &models.Articles{}
	// You can order your list here. Just change
	err := tx.All(articles)
	// to:
	// err := tx.Order("create_at desc").All(articles)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make articles available inside the html template
	c.Set("articles", articles)
	return c.Render(200, r.HTML("articles/index.html"))
}

// Show gets the data for one Article. This function is mapped to
// the path GET /articles/{article_id}
func (v ArticlesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Article
	article := &models.Article{}
	// To find the Article the parameter article_id is used.
	err := tx.Find(article, c.Param("article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make article available inside the html template
	c.Set("article", article)
	return c.Render(200, r.HTML("articles/show.html"))
}

// New renders the formular for creating a new article.
// This function is mapped to the path GET /articles/new
func (v ArticlesResource) New(c buffalo.Context) error {
	// Make article available inside the html template
	c.Set("article", &models.Article{})
	// return c.Render(200, r.HTML("articles/new.html"))
	return c.Render(200, r.HTML("articles/submit.html"))
}

// Create adds a article to the DB. This function is mapped to the
// path POST /articles
func (v ArticlesResource) Create(c buffalo.Context) error {
	// Allocate an empty Article
	article := &models.Article{}
	f := nulls.NewBool(false)
	article.Publish = f
	article.Reject = f

	// Bind article to the html form elements
	err := c.Bind(article)
	if err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(article)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make article available inside the html template
		c.Set("article", article)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		// return c.Render(422, r.HTML("articles/new.html"))
		return c.Render(422, r.HTML("articles/submit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Article was submitted successfully")
	// and redirect to the articles index page
	// return c.Redirect(302, "/articles/%s", article.ID)
	return c.Redirect(302, "/")
}

// Updated from resource template
// Edit renders a edit formular for a article. This function is
// mapped to the path GET /articles/{article_id}/edit
func (v ArticlesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Article
	article := &models.Article{}
	err := tx.Find(article, c.Param("article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make article available inside the html template
	c.Set("article", article)
	return c.Render(200, r.HTML("articles/edit.html"))
}

// Update changes a article in the DB. This function is mapped to
// the path PUT /articles/{article_id}
func (v ArticlesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Article
	article := &models.Article{}
	err := tx.Find(article, c.Param("article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind article to the html form elements
	err = c.Bind(article)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(article)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make article available inside the html template
		c.Set("article", article)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("articles/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Article was updated successfully")
	// and redirect to the articles index page
	return c.Redirect(302, "/articles/%s", article.ID)
}

// Destroy deletes a article from the DB. This function is mapped
// to the path DELETE /articles/{article_id}
func (v ArticlesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Article
	article := &models.Article{}
	// To find the Article the parameter article_id is used.
	err := tx.Find(article, c.Param("article_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(article)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Article was destroyed successfully")
	// Redirect to the articles index page
	return c.Redirect(302, "/articles")
}

// <================>Added<=================>

// New renders the formular for creating a new article.
// This function is mapped to the path GET /articles/submit
func ArticleSubmit(c buffalo.Context) error {
	// Make article available inside the html template
	c.Set("article", &models.Article{})
	// return c.Render(200, r.HTML("articles/new.html"))
	return c.Render(200, r.HTML("articles/submit.html"))
}

// List gets all Articles. This function is mapped to the path
// GET /articles
func ArticlesAdmin(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	articles := &models.Articles{}
	// You can order your list here. Just change
	// err := tx.All(articles)
	// to:
	// err := tx.Order("create_at desc").All(articles)
	err := tx.All(articles)
	if err != nil {
		return errors.WithStack(err)
	}

	// Make articles available inside the html template
	c.Set("articles", articles)
	return c.Render(200, r.HTML("articles/admin.html"))
}
