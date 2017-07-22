package actions

import (
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/i18n"

	"github.com/ECAllen/debatehub/models"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/pop"
	"github.com/pkg/errors"

	"github.com/casbin/casbin"
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
		app.Use(SetVars)

		// Casbin RBAC enforcer
		e := casbin.NewEnforcer("rbac/model.conf", "rbac/policy.csv")

		// sub := "alice" // the user that wants to access a resource.
		// obj := "data1" // the resource that is going to be accessed.
		// act := "read" // the operation that the user performs on the resource.

		if e.Enforce(sub, obj, act) == true {
			// permit alice to read data1
		} else {
			// deny the request, show an error
		}
		//---------------------
		//	Routes
		//---------------------
		app.Use(CheckLoggedIn)
		app.GET("/", func(c buffalo.Context) error {

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

		app.GET("/privacy", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("/privacy/index.md"))
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
		// TODO
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

		var ar buffalo.Resource
		ar = &ArticlesResource{&buffalo.BaseResource{}}
		// no authentication needed
		app.POST("/articles", ar.Create)
		app.GET("/articles/submit", ArticleSubmit)
		// authentication
		articles := app.Group("/articles")
		articles.Use(CheckAuth)
		articles.GET("/", ar.List)
		articles.GET("/admin", ArticlesAdmin)
		articles.GET("/new", ar.New)
		articles.GET("/{article_id}", ar.Show)
		articles.GET("/{article_id}/edit", ar.Edit)
		articles.PUT("/{article_id}", ar.Update)
		articles.DELETE("/{article_id}", ar.Destroy)

		// ------------------------
		//  Trends
		// ------------------------

		var tr buffalo.Resource
		tr = &TrendsResource{&buffalo.BaseResource{}}
		// no authentication needed
		app.GET("/trends/submit", TrendsSubmit)
		app.POST("/trends", tr.Create)
		// authentication
		trends := app.Group("/trends")
		trends.Use(CheckAuth)
		trends.GET("/", tr.List)
		trends.GET("/admin", TrendsAdmin)
		trends.GET("/new", tr.New)
		trends.GET("/{trend_id}", tr.Show)
		trends.GET("/{trend_id}/edit", tr.Edit)
		trends.PUT("/{trend_id}", tr.Update)
		trends.DELETE("/{trend_id}", tr.Destroy)

		// ------------------------
		//  Speculations
		// ------------------------

		var sp buffalo.Resource
		sp = &SpeculationsResource{&buffalo.BaseResource{}}
		// no authentication needed
		app.GET("/speculations/submit", SpeculationsSubmit)
		app.POST("/speculations", sp.Create)
		// authentication
		speculations := app.Group("/speculations")
		speculations.Use(CheckAuth)
		speculations.GET("/", sp.List)
		speculations.GET("/admin", SpeculationsAdmin)
		speculations.GET("/new", sp.New)
		speculations.GET("/{speculation_id}", sp.Show)
		speculations.GET("/{speculation_id}/edit", sp.Edit)
		speculations.PUT("/{speculation_id}", sp.Update)
		speculations.DELETE("/{speculation_id}", sp.Destroy)
	}

	return app
}
