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
		return c.Error(401, err)
	}

	// TODO remove later
	// fmt.Printf("%+v\n", user)
	// The default value just renders the data we get by GitHub
	// return c.Render(200, r.JSON(user))

	// Adding the userID to the session to remember the logged in user
	id := fmt.Sprintf("%s.%s.%s", user.Provider, user.UserID, user.Name)
	session := c.Session()
	session.Set("userID", user.UserID)

	// session.Set("userProfileId", id)
	c.Set("userProfileId", id)

	session.Set("userProvider", user.Provider)
	session.Set("userName", user.Name)
	session.Set("userNickName", user.NickName)
	session.Set("userFirstName", user.FirstName)
	session.Set("userLastName", user.LastName)
	session.Set("userAvatarURL", user.AvatarURL)

	err = session.Save()
	if err != nil {
		return c.Error(401, err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// Allocate an empty Profile
	profile := &models.Profile{}

	// Check if profile exists
	err = tx.Find(profile, id)
	if err != nil {
		return c.Redirect(http.StatusFound, "/profiles/submit")
	} else {
		// After the user is logged and profile is created then redirect
		checkPath := fmt.Sprintf("%s", session.Get("checkAuthPath"))
		return c.Redirect(http.StatusFound, checkPath)
	}
}
