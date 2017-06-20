package actions

import "github.com/gobuffalo/buffalo"

// EmailHandler for submitted emails
func EmailHandler(c buffalo.Context) error {
	c.Set("name", c.Value("emailAddr"))
	return c.Render(200, r.HTML("email.html"))
}
