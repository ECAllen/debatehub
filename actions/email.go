package actions

import "github.com/gobuffalo/buffalo"

// EmailHandler for submitted emails
func EmailHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("email.html"))
	// return c.Render(200, r.String(c.Param("emailAddr")))
}
