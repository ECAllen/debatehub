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

var forumThreadHTML = `
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
         <small><strong><a href="/profiles/{{.Profile.ID}}">{{.Profile.NickName}}</a></strong></small>
         <p class="lead">{{.Topic}}</p>
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
				threads := help.Value("threads").(Ftree)
				s := buildForum(&threads, false)
				t.Append(s)
				return t.HTML(), nil
			},
		},
	})
}
