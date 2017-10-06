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
	models.Profile
	DebateID uuid.UUID
	Token    string
	Children []Ptree
}

// TODO OK this shit is ugly... move to render.go later and refactor
var debateHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">supporting point</button>`

var debateEndHTML = `
     </div> 
     </div> <!-- close debate row -->`

var pointHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">counter point</button>`

var pointEndHTML = `
     </div> `

var counterPointHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>`

var counterPointEndHTML = `
</div> <!-- close point row -->`

var pointFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addcounterpoint?point_id={{.Point.ID}}" id="{{.Point.ID}}" method="POST" style="display:none">
	<input class="counterpoint_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<label>Topic</label>
			<textarea class=" form-control" id="debate-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
</form>`

var debateFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addpoint" id="{{.DebateID}}" method="POST" style="display:none">
	<input class="debate_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<label>Topic</label>
			<textarea class=" form-control" id="point-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
</form>`

var debateTmpl, _ = template.New("Debate").Parse(debateHTML)

var pointTmpl, _ = template.New("Point").Parse(pointHTML)

var counterPointTmpl, _ = template.New("CounterPoint").Parse(counterPointHTML)

var formTmpl, _ = template.New("Form").Parse(pointFormHTML)

var debateFormTmpl, _ = template.New("Form").Parse(debateFormHTML)

// TODO move level int ptree
func buildHTML(ptree *Ptree, level int) (string, error) {
	// TODO restructure code, oh so smelly !!!
	// TODO change level to enum
	// Slice to hold the templates and tags.
	var html []string

	// Buffer to hold template until
	// it is converted to string.
	var tpl bytes.Buffer

	switch level {
	case 1:
		html = append(html, `<div class="row"> <div class="col-md-9">`)
		debateTmpl.Execute(&tpl, ptree)
		debateFormTmpl.Execute(&tpl, ptree)
		html = append(html, tpl.String())
		html = append(html, `</div> <!-- close media body--></div> <!-- close media-->`)
		html = append(html, `</div> <!-- close column --></div> <!-- close debate row -->`)
	case 2:
		html = append(html, `<div class="row"> <div class="col-md-5">`)
		pointTmpl.Execute(&tpl, ptree)
		formTmpl.Execute(&tpl, ptree)
		html = append(html, tpl.String())
		html = append(html, `</div> <!-- close media body--></div> <!-- close media--> </div> <!-- close column--> `)
		if len(ptree.Children) == 0 {
			html = append(html, `</div> <!-- close point row -->`)
		}
	case 3:
		html = append(html, `<div class="col-md-4">`)
		counterPointTmpl.Execute(&tpl, ptree)
		formTmpl.Execute(&tpl, ptree)
		html = append(html, tpl.String())
		if len(ptree.Children) > 0 {
			for child, _ := range ptree.Children {
				counterPointTmpl.Execute(&tpl, child)
				formTmpl.Execute(&tpl, child)
				html = append(html, tpl.String())
				//TODO Left off here
			}
		}
		html = append(html, `</div> <!-- close media body--></div> <!-- close media--></div> <!-- close column -->`)
		html = append(html, `</div> <!-- close point row -->`)
	default:
		err := errors.New("should not have gotten here, check level buildHTML")
		return "", errors.WithStack(err)
	}

	// If the Point has children then recusrsivly
	// call buildHTML on them. Set couterpoint to true
	if level < 3 {
		if len(ptree.Children) > 0 {
			childLevel := level + 1
			for _, child := range ptree.Children {
				h, err := buildHTML(&child, childLevel)
				if err != nil {
					return "", errors.WithStack(err)
				}
				html = append(html, h)
			}
		}
	}

	// Join the slice into one big string
	return strings.Join(html, "\n"), nil
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

	// Collect point.
	if exists {
		err = q.First(point)
		if err != nil {
			return *point, err
		}
	}
	return *point, err
}

// TODO combine buildroot and build tree
func buildTree(id uuid.UUID, debateID uuid.UUID, tx *pop.Connection, ptree *Ptree) error {

	// Get point.
	point, err := Point(id, tx)
	if err != nil {
		return errors.WithStack(err)
	}

	// Check if there is a profile associated
	// with this point ID.
	profile2point := &models.Profiles2point{}
	q := tx.Where("point = ?", point.ID)
	exists, err := q.Exists(profile2point)
	if err != nil {
		return errors.WithStack(err)
	}

	// If there is a profile then get it
	profile := models.Profile{}

	if exists {
		err = q.First(profile2point)
		if err != nil {
			return errors.WithStack(err)
		}

		err = tx.Find(&profile, profile2point.Profile)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	ptree.Point = point
	ptree.Profile = profile
	ptree.DebateID = debateID

	// check if this point has any counterpoints
	p2c := &models.Points2counterpoint{}
	q = tx.Where("point = ?", id)
	exists, err = q.Exists(p2c)
	if err != nil {
		return errors.WithStack(err)
	}

	if exists {
		p2cs := []models.Points2counterpoint{}
		err = q.All(&p2cs)
		if err != nil {
			return errors.WithStack(err)
		}

		// Iterate through each counter point
		// create a ptree node recusrsivley call
		// buildtree on node.
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

func insertProfile2Point(pointID uuid.UUID, tx *pop.Connection, c buffalo.Context) error {

	// Associate profile with debate
	profile2point := &models.Profiles2point{}

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err := errors.New("should not have gotten here, check authentication path")
		return errors.WithStack(err)
	} else {
		profile2point.Profile = userID.(uuid.UUID)
	}

	profile2point.Point = pointID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(profile2point)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}
	return nil
}

func insertProfile2Thread(threadID uuid.UUID, tx *pop.Connection, c buffalo.Context) error {

	// Associate profile with debate
	profile2thread := &models.Profile2thread{}

	// Assume userID set otherwise should not have gotten
	// here raise error and abort
	if userID := c.Session().Get("UserID"); userID == nil {
		err := errors.New("should not have gotten here, check authentication path")
		return errors.WithStack(err)
	} else {
		profile2thread.Thread = userID.(uuid.UUID)
	}

	profile2thread.Thread = threadID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(profile2thread)
	if err != nil || verrs.HasAny() {
		return errors.WithStack(err)
	}
	return nil
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
	// TODO redo
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

	htm, err := buildHTML(&ptree, 1)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("debate_html", htm)

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

func AddThread(c buffalo.Context) error {
	// ==================================
	// Pull out params

	// The debate page id is needed in case we need to redirect
	// debate page if errors.
	debate_page_id := c.Param("debate_page_id")
	// Point_id is the existing point which this "counter point"
	// is attached.
	parent_thread_id := c.Param("parent_thread_id")

	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)

	// ==================================
	// Create the new thread.
	// ==================================
	newthread := &models.Thread{}
	newthread.Rank = 1

	// Bind thread to the html form elements
	err := c.Bind(newthread)
	if err != nil {
		return errors.WithStack(err)
	}

	// Validate the data from the html form.
	verrs, err := tx.ValidateAndCreate(newthread)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make point available inside the html template.
		c.Set("newthread", newthread)
		// Make the errors available inside the html template.
		c.Set("errors", verrs)
		// Render again the new.html template that the user can
		// correct the input.
		c.Flash().Add("warning", fmt.Sprintf("%s", verrs))
		// Redirect to the original debate page where the
		// counter point was created.
		return c.Redirect(422, "/debate_pages/%s", debate_page_id)
	}

	// If the debate ID and thread ID are the same
	// then add entry into debate2threads table
	// otherwise add entry into thread2counterthreads
	// table.
	if debate_page_id == parent_thread_id {
		debate2thread := &models.Debate2thread{}
		debate2thread.Debate, err = uuid.FromString(debate_page_id)
		if err != nil {
			return errors.WithStack(err)
		}

		// add point id
		debate2thread.Thread = newthread.ID
		verrs, err = tx.ValidateAndCreate(debate2thread)
		if err != nil {
			return errors.WithStack(err)
		}
	} else {
		// Put new thread and parent thread into
		// thread2counterthreads table.
		thread2counterthread := &models.Thread2counterthread{}

		thread2counterthread.Thread, err = uuid.FromString(parent_thread_id)
		if err != nil {
			return errors.WithStack(err)
		}

		thread2counterthread.Counterthread = newthread.ID
		verrs, err = tx.ValidateAndCreate(thread2counterthread)
		if err != nil {
			return errors.WithStack(err)
		}

		insertProfile2Thread(newthread.ID, tx, c)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// and redirect to the points index page
	return c.Redirect(302, "/debate_pages/%s", debate_page_id)
}
