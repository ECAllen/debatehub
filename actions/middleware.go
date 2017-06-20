package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
)

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

		checkPath := fmt.Sprintf("%s", session.Get("checkAuthPath"))
		fmt.Println("CheckAuth redirect path ==========>" + checkPath)

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
