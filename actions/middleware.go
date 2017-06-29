package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

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
		err := next(c)
		return err
	}
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

		// Read the userID out of the session
		userID := session.Get("userID")
		// If there is no userID redirect to the login page
		if userID == nil {
			err := c.Redirect(http.StatusTemporaryRedirect, "/login")
			return err
		}

		err = next(c)
		return err
	}
}
