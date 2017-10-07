package actions

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var r *render.Engine

// ==================================================
// Forums
// ==================================================

var forumThreadHTML = `
<ul class="media-list">
   <li class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Thread.ID}}">reply</button>`

var forumThreadEndHTML = `
     </div> <!-- close media-body -->
    </li> <!-- close media -->
   </ul> <!-- close media list -->`

var forumCounterThreadHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p>{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Thread.ID}}">reply</button>`

var forumCounterThreadEndHTML = `
     </div> <!-- close media body -->
</div><!-- close media  --> `

var forumFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addthread?parent_thread_id={{.Thread.ID}}" id="{{.Thread.ID}}" method="POST" style="display:none">
	<input class="counterthread_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<textarea class=" form-control" id="debate-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
 </form>`

var forumThreadTmpl, _ = template.New("Thread").Parse(forumThreadHTML)
var forumCounterThreadTmpl, _ = template.New("CounterThread").Parse(forumCounterThreadHTML)
var forumFormTmpl, _ = template.New("Form").Parse(forumFormHTML)

func buildForum(ftree *Ftree, counterThread bool) string {
	// Slice to hold the templates and tags.
	var html []string

	// Buffer to hold template until
	// it is converted to string.
	var tpl bytes.Buffer

	if counterThread {
		forumCounterThreadTmpl.Execute(&tpl, ftree)
	} else {
		forumThreadTmpl.Execute(&tpl, ftree)
	}

	forumFormTmpl.Execute(&tpl, ftree)

	html = append(html, tpl.String())

	// If the Thread has children then recusrsivly
	// call buildHTML on them. Set couterpoint to true
	if len(ftree.Children) > 0 {
		for _, child := range ftree.Children {
			html = append(html, buildForum(&child, true))
		}
	}

	// Append end tags to html.
	if counterThread {
		html = append(html, forumCounterThreadEndHTML)
	} else {
		html = append(html, forumThreadEndHTML)
	}

	// Join the slice into one big string
	return strings.Join(html, "\n")
}

// ==================================================
// Debates
// ==================================================

var debateHTML = `
         <p class="lead">{{.Topic}}</p>
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>`

var debateButtonHTML = `
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">supporting point</button>`

var debateFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addpoint" id="{{.DebateID}}" method="POST" style="display:none">
	<input class="debate_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<textarea class=" form-control" id="point-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
</form>`

var pointHTML = `
         <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p>{{.Topic}}</p>`

var pointButtonHTML = `
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">counter point</button>`

var counterPointHTML = `
         <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p>{{.Topic}}</p>`

var pointFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addcounterpoint?point_id={{.Point.ID}}" id="{{.Point.ID}}" method="POST" style="display:none">
	<input class="counterpoint_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<textarea class=" form-control" id="debate-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
</form>`

var pointTmpl, _ = template.New("Point").Parse(pointHTML)
var pointButtonTmpl, _ = template.New("PointButton").Parse(pointButtonHTML)
var counterPointTmpl, _ = template.New("CounterPoint").Parse(counterPointHTML)
var formTmpl, _ = template.New("Form").Parse(pointFormHTML)

var debateTmpl, _ = template.New("Debate").Parse(debateHTML)
var debateButtonTmpl, _ = template.New("DebateButton").Parse(debateButtonHTML)
var debateFormTmpl, _ = template.New("Form").Parse(debateFormHTML)

func buildHTML(ptree *Ptree, level int, userUUID uuid.UUID) (template.HTML, error) {

	// Slice to hold the templates and tags.
	var html []string

	// Buffer to hold template until
	// it is converted to string.
	var tpl bytes.Buffer

	divRowOpts := tags.Options{"class": "row row-decoration"}
	divRow := tags.New("div", divRowOpts)

	divMediaBodyOpts := tags.Options{"class": "media-body"}
	divMediaBody := tags.New("div", divMediaBodyOpts)

	// Collector div for the rows.
	divDebateOpts := tags.Options{"id": "debate"}
	divDebate := tags.New("div", divDebateOpts)

	// If point and debate are same id then this is the
	// root node and need diff html.
	if ptree.DebateID == ptree.Point.ID {
		block := tags.New("blockquote", tags.Options{})
		debateTmpl.Execute(&tpl, ptree)
		debateFormTmpl.Execute(&tpl, ptree)
		html = append(html, tpl.String())

		divOpts := tags.Options{"class": "col-md-9"}
		divDebateCol := tags.New("div", divOpts)

		block.Append(strings.Join(html, "\n"))
		divDebateCol.Append(block)

		debateButton := tags.New("button", tags.Options{"class": "btn btn-default btn-xs point-button", "value": ptree.Point.ID})
		debateButton.Append("supporting point")
		divDebateCol.Append(debateButton)

		/*
			For testing
			uid, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
			if err != nil {
				fmt.Printf("Something gone wrong: %s", err)
			}
		*/

		if userUUID == ptree.Profile.ID {
			strong := tags.New("strong", tags.Options{})
			strong.Append("edit")
			// /points/{point_id}/edit
			link := fmt.Sprintf("/debate_pages/%s/edit", ptree.DebateID)
			editLink := tags.New("a", tags.Options{"class": "edit", "href": link})
			editLink.Append(strong)
			divDebateCol.Append(editLink)
		}

		divRow.Append(divDebateCol)
		divDebate.Append(divRow)

		if len(ptree.Children) > 0 {
			childLevel := level + 1
			for _, child := range ptree.Children {
				row, err := buildHTML(&child, childLevel, userUUID)
				if err != nil {
					return row, errors.WithStack(err)
				}
				divDebate.Append(row)
			}
		}
	} else {
		// Else child of main debate.
		if level == 2 {

			divOpts := tags.Options{"class": "col-md-6"}
			divPointCol := tags.New("div", divOpts)

			// <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
			strong := tags.New("strong", tags.Options{})
			small := tags.New("small", tags.Options{})
			strong.Append(ptree.Profile.NickName)
			small.Append(strong)
			profileHref := fmt.Sprintf("/profiles/%s", ptree.Profile.ID)
			profileLink := tags.New("a", tags.Options{"class": "profile", "href": profileHref})
			profileLink.Append(small)

			// add template to col
			divMediaBody.Append(profileLink)

			if userUUID == ptree.Profile.ID {
				strong := tags.New("strong", tags.Options{})
				small := tags.New("small", tags.Options{})
				strong.Append("edit")
				small.Append(strong)
				link := fmt.Sprintf("/points/%s/edit", ptree.Point.ID)
				editLink := tags.New("a", tags.Options{"class": "edit", "href": link})
				editLink.Append(small)
				divMediaBody.Append(editLink)
			}

			p := tags.New("p", tags.Options{})
			p.Append(ptree.Topic)
			divMediaBody.Append(p)

			// execute templates and put in string
			// pointTmpl.Execute(&tpl, ptree)
			formTmpl.Execute(&tpl, ptree)
			pointButtonTmpl.Execute(&tpl, ptree)
			html = append(html, tpl.String())

			divMediaBody.Append(strings.Join(html, "\n"))
			divPointCol.Append(divMediaBody)
			divCounterPointCol := tags.New("div", tags.Options{"class": "col-md-6"})

			// call child if any
			if len(ptree.Children) > 0 {
				childLevel := level + 1
				// child := ptree.Children[0]
				// At this point level should == 3.
				for _, child := range ptree.Children {
					c, err := buildHTML(&child, childLevel, userUUID)
					if err != nil {
						return c, errors.WithStack(err)
					}
					divCounterPointCol.Append(c)
				}
			}

			// build row from two cols
			divRow.Append(divPointCol)
			divRow.Append(divCounterPointCol)
			return divRow.HTML(), nil
		} else {
			divCounterPointRowOpts := tags.Options{"class": "row"}
			divCounterPointRow := tags.New("div", divCounterPointRowOpts)

			// TODO move this to func
			// <small><strong><a class="profile" href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
			strong := tags.New("strong", tags.Options{})
			small := tags.New("small", tags.Options{})
			strong.Append(ptree.Profile.NickName)
			small.Append(strong)
			profileHref := fmt.Sprintf("/profiles/%s", ptree.Profile.ID)
			profileLink := tags.New("a", tags.Options{"class": "profile", "href": profileHref})
			profileLink.Append(small)
			divMediaBody.Append(profileLink)

			if userUUID == ptree.Profile.ID {
				strong := tags.New("strong", tags.Options{})
				small := tags.New("small", tags.Options{})
				strong.Append("edit")
				small.Append(strong)
				link := fmt.Sprintf("/points/%s/edit", ptree.Point.ID)
				editLink := tags.New("a", tags.Options{"class": "edit", "href": link})
				editLink.Append(small)
				divMediaBody.Append(editLink)

				// <a href="/points/3cae7237-32bd-413e-8240-75c181ac410e" data-method="DELETE" data-confirm="Are you sure?" class="btn btn-danger">Destroy</a>
				link = fmt.Sprintf("/points/%s", ptree.Point.ID)
				destroyLinkOpts := tags.Options{
					"class":        "edit",
					"href":         link,
					"data-method":  "DELETE",
					"data-confirm": "Are you sure?"}
				destroyLink := tags.New("a", destroyLinkOpts)
				strong = tags.New("strong", tags.Options{})
				small = tags.New("small", tags.Options{})
				small.Append("delete")
				strong.Append(small)
				destroyLink.Append(strong)
				divMediaBody.Append(destroyLink)
			}

			p := tags.New("p", tags.Options{})
			p.Append(ptree.Topic)
			divMediaBody.Append(p)

			// execute templates and put in string
			// counterPointTmpl.Execute(&tpl, ptree)
			formTmpl.Execute(&tpl, ptree)
			html = append(html, tpl.String())

			// add template to col
			divMediaBody.Append(strings.Join(html, "\n"))
			divCounterPointRow.Append(divMediaBody)
			return divCounterPointRow.HTML(), nil
		}
	}
	return divDebate.HTML(), nil
}

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),

		// Add template helpers here:
		Helpers: render.Helpers{
			"Debate": func(opts tags.Options, help plush.HelperContext) (template.HTML, error) {

				userUUID := help.Value("UserID").(uuid.UUID)
				// userID := fmt.Sprintf("%s", userUUID)

				ptree := help.Value("debate").(Ptree)
				t, err := buildHTML(&ptree, 1, userUUID)
				if err != nil {
					return t, errors.WithStack(err)
				}
				// return template.HTML(p), nil
				return t, nil
			},
			"Forum": func(opts tags.Options, help plush.HelperContext) (template.HTML, error) {
				t := tags.New("div", opts)
				threads := help.Value("threads").(Ftree)
				s := buildForum(&threads, false)
				t.Append(s)
				return t.HTML(), nil
			},
		},
	})
}
