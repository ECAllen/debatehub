package actions

import (
	"log"
	"math/rand"

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
// application is being run. Default is "production".
var ENV = envy.Get("GO_ENV", "production")
var HOST = envy.Get("GO_HOST", "http://debatehub.org")
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
			Host:        HOST,
		})

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}
		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(middleware.CSRF)
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
		app.Use(SetCurrentUser)

		//---------------------
		//	Routes
		//---------------------
		app.GET("/", func(c buffalo.Context) error {

			tagLines := []string{"Affirmative action for mental ghettos.",
				"Trumpets and SJW's need not apply.",
				"Demagogues hate this website.",
				"Elevating the new conciousness one debate at a time.",
				"An opinionated platform.",
				"Wolves not sheep.",
				"Popping filter bubbles since 2017"}

			// Set the sites motto to a random tag line.
			c.Set("motto", tagLines[rand.Intn(len(tagLines))])

			// Get the DB connection from the context
			tx := c.Value("tx").(*pop.Connection)

			// query for all published articles
			articles := &models.Articles{}
			err := tx.Where("reject = false").Where("publish = true").Order("updated_at desc").All(articles)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("articles", articles)

			// query for all published trends
			trends := &models.Trends{}
			err = tx.Where("reject = false").Where("publish = true").All(trends)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("trends", trends)

			// query for all published speculations
			speculations := &models.Speculations{}
			err = tx.Where("reject = false").Where("publish = true").All(speculations)
			if err != nil {
				return errors.WithStack(err)
			}
			c.Set("speculations", speculations)

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

		app.GET("/robots.txt", func(c buffalo.Context) error {
			return c.Render(200, r.Plain("/robots.txt"))
		})

		// Returns an error stack just to print out useful info.
		app.GET("/context", func(c buffalo.Context) error {
			err := errors.New("Context")
			return errors.WithStack(err)
		})

		// -----------------
		//   Authentication
		// -----------------
		auth := app.Group("/auth")
		authHandler := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", authHandler)
		auth.GET("/{provider}/callback", AuthCallback)
		app.GET("/login",
			func(c buffalo.Context) error {
				return c.Render(200, r.HTML("login/index.html"))
			})
		app.POST("/login",
			func(c buffalo.Context) error {
				return c.Render(200, r.HTML("login/index.html"))
			})
		app.DELETE("/", AuthDestroy)

		//---------------------
		//	Authorization
		//---------------------

		// ------------------
		//   Secure Content
		// ------------------

		// ------------------
		//  Profiles
		// ------------------
		app.GET("/profiles/submit", ProfilesSubmit)
		app.GET("/profiles/user", ProfileUserShow)
		var pr buffalo.Resource
		pr = &ProfilesResource{&buffalo.BaseResource{}}
		profiles := app.Resource("/profiles", pr)
		profiles.Use(CheckAuth, CheckAdmin)
		profiles.Middleware.Skip(CheckAdmin, pr.Create, pr.Show, pr.Update, pr.Edit)

		// ------------------------
		//   Email Subscriptions
		// ------------------------

		var er buffalo.Resource
		er = &EmailsResource{&buffalo.BaseResource{}}
		subscription := app.Resource("/emails", er)
		subscription.Use(CheckAuth)
		subscription.Middleware.Skip(CheckAuth, er.Create)
		subscription.Middleware.Skip(CheckAdmin, er.Create)
		// TODO fix this
		// subscription.Middleware.Skip(middleware.CSRF, er.Create)

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
		articles.Use(CheckAdmin)
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
		trends.Use(CheckAdmin)
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
		speculations.Use(CheckAdmin)
		speculations.GET("/", sp.List)
		speculations.GET("/admin", SpeculationsAdmin)
		speculations.GET("/new", sp.New)
		speculations.GET("/{speculation_id}", sp.Show)
		speculations.GET("/{speculation_id}/edit", sp.Edit)
		speculations.PUT("/{speculation_id}", sp.Update)
		speculations.DELETE("/{speculation_id}", sp.Destroy)

		// ------------------------
		//  Debates
		// ------------------------

		dbt := &DebatesResource{&buffalo.BaseResource{}}
		debates := app.Resource("/debates", dbt)
		debates.Use(CheckAuth, CheckAdmin)

		var db buffalo.Resource
		db = &DebatePagesResource{&buffalo.BaseResource{}}
		debate_pages := app.Group("/debate_pages")
		debate_pages.Use(CheckAuth)
		debate_pages.Middleware.Skip(CheckAuth, db.List, db.Show)
		debate_pages.GET("/", db.List)
		debate_pages.POST("/", db.Create)
		debate_pages.GET("/new", db.New)
		debate_pages.GET("/{debate_page_id}", db.Show)
		debate_pages.GET("/{debate_page_id}/edit", db.Edit)
		debate_pages.POST("/{debate_page_id}/addpoint", AddPoint)
		debate_pages.POST("/{debate_page_id}/addcounterpoint", AddCounterPoint)
		debate_pages.POST("/{debate_page_id}/addthread", AddThread)
		debate_pages.PUT("/{debate_page_id}", db.Update)
		debate_pages.DELETE("/{debate_page_id}", db.Destroy)

		var pt buffalo.Resource
		pt = &PointsResource{&buffalo.BaseResource{}}
		points := app.Resource("/points", pt)
		points.Use(CheckAuth)
		// profiles.Middleware.Skip(CheckAdmin, pr.Create, pr.Show, pr.Update, pr.Edit)
	}
	return app
}
