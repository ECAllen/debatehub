package actions

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

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

// Points tree for debates
type Ptree struct {
	models.Point
	DebateID    uuid.UUID
	ProfileID   uuid.UUID
	ProfileNick string
	Token       string
	Children    []Ptree
}

// TODO OK this shit is ugly... move to render.go later and refactor
var pointHTML = `
   <li class="media"> 
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <h4 class="media-heading">{{.ProfileNick}}</h4>
         <p>{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.ID}}">reply</button>`

var pointEndHTML = `
     </div> 
    </li>`

var counterPointHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <h4 class="media-heading">{{.ProfileNick}}</h4>
         <p>{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.ID}}">reply</button>`

var counterPointEndHTML = `
     </div>
</div>`

var formHTML = `
<form action="/debate_pages/{{.DebateID}}/addcounterpoint?point_id={{.ID}}" id="{{.ID}}" method="POST" style="display:none">
	<input class="counterpoint_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<label>Topic</label>
			<textarea class=" form-control" id="debate-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
 </form>`

var debateFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addpoint" id="{{.ID}}" method="POST" style="display:none">
	<input class="debate_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<label>Topic</label>
			<textarea class=" form-control" id="point-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
</form>`

var pointTmpl, _ = template.New("Point").Parse(pointHTML)

var counterPointTmpl, _ = template.New("CounterPoint").Parse(counterPointHTML)

var formTmpl, _ = template.New("Form").Parse(formHTML)

var debateFormTmpl, _ = template.New("Form").Parse(debateFormHTML)

func buildHTML(ptree *Ptree, counterPoint bool) string {
	// Slice to hold the templates and tags.
	var html []string

	// Buffer to hold template until
	// it is converted to string.
	var tpl bytes.Buffer

	// Bind vars in template for beginning.
	if counterPoint {
		counterPointTmpl.Execute(&tpl, ptree)
	} else {
		pointTmpl.Execute(&tpl, ptree)
	}

	// Build form and append
	if counterPoint {
		formTmpl.Execute(&tpl, ptree)
	} else {
		debateFormTmpl.Execute(&tpl, ptree)
	}

	html = append(html, tpl.String())

	// If the Point has children then recusrsivly
	// call buildHTML on them. Set couterpoint to true
	if len(ptree.Children) > 0 {
		for _, child := range ptree.Children {
			html = append(html, buildHTML(&child, true))
		}
	}

	// Append end tags to html.
	if counterPoint {
		html = append(html, counterPointEndHTML)
	} else {
		html = append(html, pointEndHTML)
	}

	// Join the slice into one big string
	return strings.Join(html, "\n")
}

func Point(id uuid.UUID, tx *pop.Connection) (models.Point, error) {

	// point used to hold the point
	point := &models.Point{}

	// Create query.
	q := tx.Where("ID = ?", id)

	// verify that the point exists in
	// the Point table
	exists, err := q.Exists(point)
	if err != nil {
		return *point, err
	}

	// Collect point into
	// points slice.
	if exists {
		err = q.First(point)
		if err != nil {
			return *point, err
		}
	}
	return *point, err
}

func buildTree(id uuid.UUID, debateID uuid.UUID, tx *pop.Connection, ptree *Ptree) error {

	// Get point.
	p, err := Point(id, tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// Put point into ptree.
	ptree.Point = p
	ptree.DebateID = debateID

	// check if this point has any counterpoints
	p2c := &models.Points2counterpoint{}
	q := tx.Where("point = ?", id)
	exists, err := q.Exists(p2c)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		// go through each counter point
		// create a ptree node
		// buildtree on node
		p2cs := []models.Points2counterpoint{}
		err = q.All(&p2cs)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, c := range p2cs {
			var pt Ptree
			pt.Token = ptree.Token
			err = buildTree(c.Counterpoint, debateID, tx, &pt)
			if err != nil {
				return errors.WithStack(err)
			}
			ptree.Children = append(ptree.Children, pt)
		}

	}
	return errors.WithStack(err)
}

func buildTreeRoot(id uuid.UUID, tx *pop.Connection, ptree *Ptree) error {

	// Check if this debate has any points.
	debatePoint := &models.Debates2point{}
	q := tx.Where("debate = ?", id)
	exists, err := q.Exists(debatePoint)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		// The debatePoints used to iterate through
		// to collect all the debate point uuid's.
		debatePoints := []models.Debates2point{}

		// Collect all the point id's associated
		// with the debate.
		err = q.All(&debatePoints)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, dp := range debatePoints {
			var pt Ptree
			pt.Token = ptree.Token
			err = buildTree(dp.Point, id, tx, &pt)
			if err != nil {
				return errors.WithStack(err)
			}
			ptree.Children = append(ptree.Children, pt)
		}
	}
	return errors.WithStack(err)
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

	profile2debate := []models.Profiles2debate{}
	q := tx.Where("debate = ?", debate_id)
	err = q.All(&profile2debate)
	if err != nil {
		return errors.WithStack(err)
	}
	// ==================================
	// Query for the profile
	// ==================================

	profile := &models.Profile{}
	err = tx.Find(profile, profile2debate[0].Profile)
	if err != nil {
		return errors.WithStack(err)
	}

	// ==================================
	// Create the root node of Ptree
	// ==================================

	var ptree Ptree
	ptree.Topic = debate.Topic
	ptree.ID = debate.ID
	ptree.DebateID = debate.ID
	ptree.ProfileID = profile.ID
	ptree.ProfileNick = profile.NickName
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

	htm := buildHTML(&ptree, false)

	c.Set("debate", debate)
	c.Set("debate_html", htm)

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

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err = errors.New("should not have gotten here, check authentication")
		return errors.WithStack(err)
	} else {
		profile2debate.Profile = userID.(uuid.UUID)
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

	profile2debate.Debate = debate.ID

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
	// fmt.Println("point ================" + fmt.Sprintf("%s", point))
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

	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_page_id)
}
