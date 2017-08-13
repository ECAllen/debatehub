package actions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		twitter.New(
			os.Getenv("TWITTER_KEY"),
			os.Getenv("TWITTER_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/twitter/callback")),
		discord.New(
			os.Getenv("DISCORD_KEY"),
			os.Getenv("DISCORD_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/discord/callback")),
		github.New(
			os.Getenv("GITHUB_KEY"),
			os.Getenv("GITHUB_SECRET"),
			fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())

	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// TODO remove later
	// fmt.Printf("%+v\n", user)
	// The default value just renders the data we get by GitHub
	// return c.Render(200, r.JSON(user))

	q := tx.Where("provider = ? and userid = ?", user.Provider, user.UserID)
	exists, err := q.Exists("profiles")

	if err != nil {
		return errors.WithStack(err)
	}

	// Allocate an empty Profile
	profile := &models.Profile{}

	if exists {
		err := q.First(profile)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		profile.UserID = user.UserID
		profile.Provider = user.Provider
		profile.NickName = user.NickName
		profile.FirstName = user.FirstName
		profile.LastName = user.LastName
		profile.Email = user.Email
		profile.AvatarURL = nulls.NewString(user.AvatarURL)
	}

	err = tx.Save(profile)
	if err != nil {
		return errors.WithStack(err)
	}

	// Adding the user info to the session
	session := c.Session()
	session.Set("UserID", user.UserID)
	err = session.Save()
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		checkPath := fmt.Sprintf("%s", session.Get("checkAuthPath"))
		return c.Redirect(http.StatusFound, checkPath)
	} else {
		c.Flash().Add("success", "Please create your profile.")
		return c.Redirect(http.StatusFound, "/profiles/submit")
	}
}

// SetCurrentUser finds and sets the logged in user
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("UserID"); uid != nil {
			p := &models.Profile{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(p, uid)
			if err != nil {
				errors.WithStack(err)
			}
			c.Set("CurrentUser", p)
			c.Set("UserID", p.ID)
			c.Set("loggedin", true)
		} else {
			c.Set("loggedin", false)
		}
		return next(c)
	}
}

// Destroy the session data upon logout
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	err := c.Session().Save()
	if err != nil {
		return errors.WithStack(err)
	}
	c.Flash().Add("success", "You have been logged out")
	return c.Redirect(http.StatusFound, "/")
}

// CheckAuth is the middleware to check if a user is logged on.
func CheckAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		// Store current path in the context so that AuthCallback
		// can redirect to correct url once authenticated
		p := fmt.Sprintf("%s", c.Value("current_path"))
		session := c.Session()
		session.Set("checkAuthPath", p)
		err := session.Save()
		if err != nil {
			return c.Error(401, err)
		}

		// If there is no userID redirect to the login page
		if userID := session.Get("UserID"); userID == nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		return next(c)
	}
}

func CheckLoggedIn(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Read the userID out of the session
		userID := c.Session().Get("userID")
		// If there is no userID redirect to the login page
		if userID == nil {
			c.Set("loggedin", false)
		} else {
			c.Set("loggedin", true)
		}
		return next(c)
	}
}
