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

		app.Use(T.Middleware())
		app.GET("/", HomeHandler)

		app.GET("/blog/{post}", func(c buffalo.Context) error {
			return c.Render(200, r.HTML("blog/"+c.Param("post")+".md"))
		})

		// app.Resource("/users", UsersResource{&buffalo.BaseResource{}})
		app.Resource("/emails", EmailsResource{&buffalo.BaseResource{}})
		app.ServeFiles("/assets", packr.NewBox("../public/assets"))

		auth := app.Group("/auth")
		auth.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
		auth.GET("/{provider}/callback", AuthCallback)
	}

	return app
}
