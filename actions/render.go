package actions

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
)

var r *render.Engine

var forumPointHTML = `
<ul class="media-list">
   <li class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">reply</button>`

var forumPointEndHTML = `
     </div> <!-- close media-body -->
    </li> <!-- close media -->
   </ul> <!-- close media list -->`

var forumCounterPointHTML = `
<div class="media">
     <div class="media-left">
       <a href="#">
         <img class="media-object" src="" alt="">
       </a>
     </div>
     <div class="media-body">
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
	 <button class="btn btn-default btn-xs point-button" value="{{.Point.ID}}">reply</button>`

var forumCounterPointEndHTML = `
     </div> <!-- close media body -->
</div><!-- close media  --> `

var forumFormHTML = `
<form action="/debate_pages/{{.DebateID}}/addcounterpoint?point_id={{.Point.ID}}" id="{{.Point.ID}}" method="POST" style="display:none">
	<input class="counterpoint_form" name="authenticity_token" value="{{.Token}}" type="hidden">
   		<div class="form-group">
			<label>Reply</label>
			<textarea class=" form-control" id="debate-Topic" name="Topic" rows="3"></textarea>
		</div>
		<button class="btn btn-success" role="submit">Add</button>
 </form>`

var forumPointTmpl, _ = template.New("Point").Parse(forumPointHTML)

var forumCounterPointTmpl, _ = template.New("CounterPoint").Parse(forumCounterPointHTML)

var forumFormTmpl, _ = template.New("Form").Parse(forumFormHTML)

func buildForum(ptree *Ptree, counterPoint bool) string {
	// Slice to hold the templates and tags.
	var html []string

	// Buffer to hold template until
	// it is converted to string.
	var tpl bytes.Buffer

	if counterPoint {
		forumCounterPointTmpl.Execute(&tpl, ptree)
	} else {
		forumPointTmpl.Execute(&tpl, ptree)
	}

	forumFormTmpl.Execute(&tpl, ptree)

	html = append(html, tpl.String())

	// If the Point has children then recusrsivly
	// call buildHTML on them. Set couterpoint to true
	if len(ptree.Children) > 0 {
		for _, child := range ptree.Children {
			html = append(html, buildForum(&child, true))
		}
	}

	// Append end tags to html.
	if counterPoint {
		html = append(html, forumCounterPointEndHTML)
	} else {
		html = append(html, forumPointEndHTML)
	}

	// Join the slice into one big string
	return strings.Join(html, "\n")
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
				p := help.Value("debate_html").(string)
				return template.HTML(p), nil
			},
			"Forum": func(opts tags.Options, help plush.HelperContext) (template.HTML, error) {
				t := tags.New("div", opts)
				threads := help.Value("threads").(Ptree)
				s := buildForum(&threads, false)
				t.Append(s)
				return t.HTML(), nil
			},
		},
	})
}
