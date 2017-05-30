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
		})

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		// app.Use(middleware.CSRF)

		app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		T, err = i18n.New(packr.NewBox("../locales"), "en-US")
		if err != nil {
			log.Fatal(err)
		}

		// TODO review all URL paths for authorization, use a grift
		app.ServeFiles("/assets", packr.NewBox("../public/assets"))
		app.Use(T.Middleware())
		app.GET("/", HomeHandler)

		app.GET("/blog/{post}", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("blog/"+c.Param("post")+".md"))
		})

		// TODO move users to auth
		// app.Resource("/users", UsersResource{&buffalo.BaseResource{}})

		// TODO verify emails URL paths needs auth
		app.Resource("/emails", EmailsResource{&buffalo.BaseResource{}})

		// -----------------
		//   Authorization
		// -----------------
		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)
		auth.GET("/login",
			func(c buffalo.Context) error {
				return c.Render(200, r.HTML("login/index.html"))
			})

		// ------------------
		//   Secure Content
		// ------------------
		secure := app.Group("/secure")
		secure.Use(CheckAuth)
		secure.GET("/",
			func(c buffalo.Context) error {
				return c.Render(200, r.HTML("secure/index.html"))
			})

		secure.DELETE("/logout",
			func(c buffalo.Context) error {
				session := c.Session()
				session.Delete("userID")
				session.Save()
				return c.Redirect(301, "/login")
			})

		app.Redirect(301, "/", "/login")
	}

	return app
}
