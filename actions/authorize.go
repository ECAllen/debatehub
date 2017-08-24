package actions

import (
	"github.com/gobuffalo/buffalo"
	"net/http"
)

func CheckAdmin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {

		// Check if user logged in and if so then
		if nick := c.Value("Nick"); nick == nil || nick != "therewillbewar" {
			c.Flash().Add("success", "You do not have permissions to view this page.")
			return c.Redirect(http.StatusFound, "/")
		}

		return next(c)
	}
}
