package actions

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/i18n"

	"github.com/ECAllen/debatehub/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// T i18n translator see locales/
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:         ENV,
			SessionName: "_debatehub_session",
			// Host:        "http://localhost:3000",
		})

		// TODO update ENV for deployment
		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		// TODO
		// app.Use(middleware.CSRF)

		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		T, err = i18n.New(packr.NewBox("../locales"), "en-US")
		if err != nil {
			log.Fatal(err)
		}

		app.ServeFiles("/assets", packr.NewBox("../public/assets"))
		app.Use(T.Middleware())

		//---------------------
		//	Routes
		//---------------------
		app.Use(CheckLoggedIn)
		app.GET("/", func(c buffalo.Context) error {

			year := fmt.Sprintf("%d", time.Now().Year())
			c.Set("year", year)

			// Get the DB connection from the context
			tx := c.Value("tx").(*pop.Connection)
			articles := &models.Articles{}
			trends := &models.Trends{}

			// query for all published articles
			err := tx.Where("reject = false").Where("publish = true").All(articles)
			if err != nil {
				return errors.WithStack(err)
			}

			// query for all published trends
			err = tx.Where("reject = false").Where("publish = true").All(trends)
			if err != nil {
				return errors.WithStack(err)
			}

			// Make articles available inside the html template
			c.Set("articles", articles)
			c.Set("trends", trends)

			return c.Render(200, r.HTML("index.html"))
		})

		app.GET("/blog", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("/blog/index.md"))
		})

		// TODO check this for injection
		app.GET("/blog/{post}", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("blog/"+c.Param("post")+".md"))
		})

		app.GET("/mission", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("/mission/index.md"))
		})

		// -----------------
		//   Authentication
		// -----------------
		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)
		app.GET("/login",
			func(c buffalo.Context) error {
				return c.Render(200, r.HTML("login/index.html"))
			})
		app.DELETE("/logout",
			func(c buffalo.Context) error {
				session := c.Session()
				session.Delete("userID")
				session.Save()
				return c.Redirect(301, "/")
			})

		// ------------------
		//   Secure Content
		// ------------------

		// ------------------
		//  Profiles
		// ------------------

		profiles := app.Resource("/profiles", ProfilesResource{&buffalo.BaseResource{}})
		profiles.Use(CheckAuth)

		// ------------------------
		//   Email Subscriptions
		// ------------------------

		var er buffalo.Resource
		er = &EmailsResource{&buffalo.BaseResource{}}
		subscription := app.Resource("/emails", er)
		subscription.Use(CheckAuth)
		subscription.Middleware.Skip(CheckAuth, er.Create)

		// ------------------------
		//  Articles
		// ------------------------

		app.GET("/articles/submit", ArticleSubmit)
		app.GET("/articles/admin", ArticlesAdmin)
		var ar buffalo.Resource
		ar = &ArticlesResource{&buffalo.BaseResource{}}
		articles := app.Resource("/articles", ar)
		articles.Use(CheckAuth)
		articles.Middleware.Skip(CheckAuth, ar.Create)

		// ------------------------
		//  Trends
		// ------------------------

		app.GET("/trends/submit", TrendsSubmit)
		app.GET("/trends/admin", TrendsAdmin)
		var tr buffalo.Resource
		tr = &TrendsResource{&buffalo.BaseResource{}}
		trends := app.Resource("/trends", tr)
		trends.Use(CheckAuth)
		trends.Middleware.Skip(CheckAuth, tr.Create)
	}

	return app
}
