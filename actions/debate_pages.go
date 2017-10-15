package actions

import (
	"fmt"

	"github.com/ECAllen/debatehub/models"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// DebatePagesResource is the resource for the debate model
type DebatePagesResource struct {
	buffalo.Resource
}

// List gets all Debates. This function is mapped to the path
// GET /debate_pages
func (v DebatePagesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	debates := &models.Debates{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// You can order your list here. Just change
	err := q.All(debates)
	// to:
	// err := q.Order("created_at desc").All(debate_pages)
	if err != nil {
		return errors.WithStack(err)
	}
	// Make Debates available inside the html template
	c.Set("debates", debates)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("debate_pages/index.html"))
}

// Show gets the data for one Debate. This function is mapped to
// the path GET /debate_pages/{debate_page_id}
func (v DebatePagesResource) Show(c buffalo.Context) error {
	// ==================================
	// Setup
	// ==================================

	// Get params
	debate_id := c.Param("debate_page_id")
	tx := c.Value("tx").(*pop.Connection)
	auth_token := c.Value("authenticity_token").(string)

	// ==================================
	// Query for the debate.
	// ==================================

	debate := &models.Debate{}
	err := tx.Find(debate, debate_id)
	if err != nil {
		return errors.WithStack(err)
	}

	// ==================================
	// Query for the profile ID for debate ID
	// ==================================

	profile2debate := models.Profiles2debate{}
	q := tx.Where("debate = ?", debate_id)
	err = q.First(&profile2debate)
	if err != nil {
		return errors.WithStack(err)
	}
	// ==================================
	// Query for the profile
	// ==================================

	profile := &models.Profile{}
	err = tx.Find(profile, profile2debate.Profile)
	if err != nil {
		return errors.WithStack(err)
	}

	// ==================================
	// Create the root node of Ptree
	// ==================================

	var ptree Ptree
	ptree.Topic = debate.Topic
	ptree.Point.ID = debate.ID
	ptree.DebateID = debate.ID
	ptree.Profile = *profile
	ptree.Token = auth_token

	// ==================================
	// Counter points
	// ==================================

	// Check for the existence of counter
	// points for the debate in the debates2point
	// table.
	debates2point := &models.Debates2point{}
	q = tx.Where("debate = ?", debate_id)
	exists, err := q.Exists(debates2point)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there are counter points then build
	// points tree.
	if exists {
		// Build the tree of comments.
		err = buildTreeRoot(debate.ID, tx, &ptree)
		if err != nil {
			return err
		}
	}
	/*
		htm, err := buildHTML(&ptree, 1)
		if err != nil {
			return errors.WithStack(err)
		}
	*/
	c.Set("debate", ptree)

	// ==================================
	// Create the root node of Ftree
	// ==================================

	var ftree Ftree
	ftree.Topic = debate.Topic
	ftree.Thread.ID = debate.ID
	ftree.DebateID = debate.ID
	ftree.Profile = *profile
	ftree.Token = auth_token

	// check for counter threads

	// Check for the existence of counter
	// threads for the forum in the debates2threads
	// table.
	debates2thread := &models.Debate2thread{}
	q = tx.Where("debate = ?", debate_id)
	exists, err = q.Exists(debates2thread)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there are counter threads then build
	// Ftree.
	if exists {
		// Build the tree of comments.
		err = buildFTreeRoot(debate.ID, tx, &ftree)
		if err != nil {
			return err
		}
	}

	c.Set("threads", ftree)

	return c.Render(200, r.HTML("debate_pages/show.html"))
}

// New renders the formular for creating a new Debate.
// This function is mapped to the path GET /debate_pages/new
func (v DebatePagesResource) New(c buffalo.Context) error {
	// Make debate available inside the html template
	c.Set("debate", &models.Debate{})
	return c.Render(200, r.HTML("debate_pages/new.html"))
}

// Create adds a Debate to the DB. This function is mapped to the
// path POST /debate_pages
func (v DebatePagesResource) Create(c buffalo.Context) error {
	// Allocate an empty Debate
	debate := &models.Debate{}
	profile2debate := &models.Profiles2debate{}

	// Bind debate to the html form elements
	err := c.Bind(debate)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(debate)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make debate available inside the html template
		c.Set("debate", debate)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("debate_pages/new.html"))
	}

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err = errors.New("should not have gotten here, check authentication")
		return errors.WithStack(err)
	} else {
		profile2debate.Profile = userID.(uuid.UUID)
	}

	profile2debate.Debate = debate.ID

	fmt.Printf("\n\n%v\n", debate.ID)

	// Validate the data from the html form
	verrs, err = tx.ValidateAndCreate(profile2debate)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Debate was created successfully")

	// and redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages/%s", debate.ID)
}

// Edit renders a edit formular for a debate. This function is
// mapped to the path GET /debate_pages/{debate_page_id}/edit
func (v DebatePagesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Make debate available inside the html template
	c.Set("debate", debate)
	return c.Render(200, r.HTML("debate_pages/edit.html"))
}

// Update changes a debate in the DB. This function is mapped to
// the path PUT /debate_pages/{debate_page_id}
func (v DebatePagesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Debate to the html form elements
	err = c.Bind(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make debate available inside the html template
		c.Set("debate", debate)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("debate_pages/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Debate was updated successfully")
	// and redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages/%s", debate.ID)
}

// Destroy deletes a debate from the DB. This function is mapped
// to the path DELETE /debate_pages/{debate_page_id}
func (v DebatePagesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Debate
	debate := &models.Debate{}
	// To find the Debate the parameter debate_page_id is used.
	err := tx.Find(debate, c.Param("debate_page_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(debate)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Debate was destroyed successfully")
	// Redirect to the debate_pages index page
	return c.Redirect(302, "/debate_pages")
}

func AddPoint(c buffalo.Context) error {

	debate_id := c.Param("debate_page_id")

	// Allocate an empty Point
	point := &models.Point{}
	point.Rank = 1

	// Bind point to the html form elements
	err := c.Bind(point)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(point)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make point available inside the html template
		c.Set("point", point)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		// Adding the user info to the session
		fmt.Println(verrs)
		c.Flash().Add("warning", fmt.Sprintf("%s", verrs))
		return c.Redirect(422, "/debate_pages/%s", debate_id)
	}

	// put point and debate into points2counterpoints table
	// add debate id
	debate2point := &models.Debates2point{}
	debate2point.Debate, err = uuid.FromString(debate_id)
	if err != nil {
		return errors.WithStack(err)
	}

	// add point id
	debate2point.Point = point.ID
	verrs, err = tx.ValidateAndCreate(debate2point)
	if err != nil {
		return errors.WithStack(err)
	}

	insertProfile2Point(point.ID, tx, c)
	if err != nil {
		return errors.WithStack(err)
	}

	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_id)
}

func AddCounterPoint(c buffalo.Context) error {
	// ==================================
	// Pull out params

	// The debate page id is needed in case we need to redirect
	// debate page if errors.
	debate_page_id := c.Param("debate_page_id")
	// Point_id is the existing point which this "counter point"
	// is attached.
	point_id := c.Param("point_id")

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// ==================================
	// Create the counter point.

	counterpoint := &models.Point{}
	counterpoint.Rank = 1
	// Bind point to the html form elements
	err := c.Bind(counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}

	// Validate the data from the html form.
	verrs, err := tx.ValidateAndCreate(counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make point available inside the html template.
		c.Set("counterpoint", counterpoint)
		// Make the errors available inside the html template.
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		c.Flash().Add("warning", fmt.Sprintf("%s", verrs))
		// Redirect to the original debate page where the
		// counter point was created.
		return c.Redirect(422, "/debate_pages/%s", debate_page_id)
	}

	// ==================================
	// Put point and counter point into points2counterpoints table.
	point2counterpoint := &models.Points2counterpoint{}

	point2counterpoint.Point, err = uuid.FromString(point_id)
	if err != nil {
		return errors.WithStack(err)
	}

	// add point id
	point2counterpoint.Counterpoint = counterpoint.ID
	verrs, err = tx.ValidateAndCreate(point2counterpoint)
	if err != nil {
		return errors.WithStack(err)
	}

	insertProfile2Point(counterpoint.ID, tx, c)
	if err != nil {
		return errors.WithStack(err)
	}

	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_page_id)
}

// Edit renders a edit formular for a point. This function is
// mapped to the path GET /debate_pages/{point_id}/pointedit?debate_page_id=
func PointEdit(c buffalo.Context) error {

	// The debate page id is needed in case we need to redirect
	// debate page if errors.
	debate_page_id := c.Param("debate_page_id")

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Point
	point := &models.Point{}
	err := tx.Find(point, c.Param("point_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	// Make point available inside the html template
	action := fmt.Sprintf("/debate_pages/%s/pointupdate?debate_page_id=%s", point.ID, debate_page_id)
	c.Set("point", point)
	c.Set("action", action)
	return c.Render(200, r.HTML("debate_pages/pointedit.html"))
}

// Update changes a point in the DB. This function is mapped to
// the path PUT /debate_pages/{point_id}/updatepoint?debate_page_id=
func PointUpdate(c buffalo.Context) error {

	// The debate page id is needed in case we need to redirect
	// debate page if errors.
	debate_page_id := c.Param("debate_page_id")

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Point
	point := &models.Point{}
	err := tx.Find(point, c.Param("point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	// Bind Point to the html form elements
	err = c.Bind(point)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := tx.ValidateAndUpdate(point)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// Make point available inside the html template
		c.Set("point", point)
		// Make the errors available inside the html template
		c.Set("errors", verrs)
		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.HTML("points/edit.html"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "Point was updated successfully")
	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_page_id)
}

// Destroy deletes a point from the DB. This function is mapped
// to the path DELETE /debate_pages/{point_id}/pointdestroy?debate_page_id=
func PointDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Allocate an empty Point
	point := &models.Point{}
	// To find the Point the parameter point_id is used.
	err := tx.Find(point, c.Param("point_id"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = tx.Destroy(point)
	if err != nil {
		return errors.WithStack(err)
	}
	// If there are no errors set a flash message
	c.Flash().Add("success", "Point was destroyed successfully")
	// Redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", c.Param("debate_page_id"))
}

func Article(c buffalo.Context) error {
	// get param
	title := c.Param("title")

	// lookup up debate by article title
	debate := &models.Debate{}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("title = ?", title)
	exists, err := q.Exists(debate)
	if err != nil {
		return errors.WithStack(err)

	}

	// If exists then route to debate id.
	if exists {
		// Then lookup debate
		err = q.First(debate)
		if err != nil {
			return errors.WithStack(err)
		}
		return c.Redirect(302, "/debate_pages/%s", debate.ID)
	}

	// Make debate available inside the html template
	debate.Title = title
	debate.Topic = title
	c.Set("debate", debate)
	return c.Render(200, r.HTML("debate_pages/newArticle.html"))
}
