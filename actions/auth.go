package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/twitter"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/twitter/callback")),
		discord.New(os.Getenv("DISCORD_KEY"), os.Getenv("DISCORD_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/discord/callback")),
	)
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}
