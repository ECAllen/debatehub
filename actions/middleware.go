package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

// CheckAuth is the middlewar to check if a user is logged on.
func CheckAuth(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// Read the userID out of the session
		userID := c.Session().Get("userID")
		// If there is no userID redirect to the login page
		if userID == nil {
			err := c.Redirect(http.StatusTemporaryRedirect, "/login")
			return err
		}
		err := next(c)
		return err
	}
}
