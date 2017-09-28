package actions

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/tags"
	"html/template"
)

var r *render.Engine

/*
 25 <ul class="media-list">
 24 <%= for (point) in points { %>
 23 <div class="panel panel-default">
 22   <li class="media">
 21     <div class="media-left">
 20       <a href="#">
 19         <img class="media-object" src="" alt="">
 18       </a>
 17     </div>
 16     <div class="media-body">
 15             <p class="media-heading"><%= point.Rank %></p>
 14         <!-- insert here for sub -->
 13         <p><%= point.Topic %></p>


 12 <%= form_for(debate, {action: debatePageAddcounterpointPath({ debate_page_id: debate.ID, point_id: p
    oint.ID }), method: "POST"}) { %>
 11   <%= partial("debate_pages/pointForm.html") %>
 10   <a href="<%= debatePageAddcounterpointPath() %>" class="btn btn-warning" data-confirm="Are you sur
    e?">Cancel</a>
  9 <% } %>


  8         </div>
  7     </div>
  6   </li>
  5     </div>
  4 <% } %>
  3 </ul>
*/

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),

		// Add template helpers here:
		Helpers: render.Helpers{
			"buildDebate": func(opts tags.Options, help plush.HelperContext) (template.HTML, error) {
				p := help.Value("debate_html").(string)
				return template.HTML(p), nil
			},
		},
	})
}
