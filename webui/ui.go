package webui

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/gregoryv/miniplan"
	. "github.com/gregoryv/miniplan"
)

func NewUI(sys *System) *UI {
	http.HandleFunc("/static/theme.css", serveTheme)
	http.HandleFunc("/static/tools.js", serveTools)

	ui := &UI{System: sys}
	http.Handle("/", ui)
	return ui
}

type UI struct {
	*miniplan.System
}

func (me *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, _ := me.Query("SELECT * FROM changes")

		defer rows.Close()
		var changes []Change
		for rows.Next() {
			var c Change
			rows.Scan(&c.UUID, &c.Title, &c.Description)
			changes = append(changes, c)
		}
		m := map[string]interface{}{
			"Changes": changes,
		}
		err := tpl.Execute(w, m)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

	case "POST":
		switch r.PostFormValue("submit") {
		case "add":
			c := Change{
				Title:       r.PostFormValue("title"),
				Description: r.PostFormValue("description"),
			}
			_ = me.Create(&c)

		case "delete":
			err := me.Remove(r.PostFormValue("uuid"))
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
		}
		http.Redirect(w, r, "/", 303)

	default:
		w.WriteHeader(405)
	}
}

var index = `
<!doctype html>

<html>
<head><title>miniplan</title>

<link rel="stylesheet" type="text/css" href="/static/theme.css" />
<script src="/static/tools.js"></script>

</head>
<body id="body">

<h1>Miniplan</h1>

{{range .Changes}}

<div class="row entry">

<div class="left">
<a href="#{{.Ref}}" class="idref">#</a>
<a name="{{.Ref}}">{{.Ref}}</a>
</div>

<div class="mid">
{{.Title}}<br>
<p>{{.Description}}</p>
</div>


<div class="right">
<form method="POST" class="one">
<input type=hidden name=submit value=insert>
<input type=submit value=I>
</form>
<form method="POST" class="one">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=hidden name=submit value=delete>
<input type=submit value=D>
</form>
</div>

</div>
{{end}}


<div class="row">
<div class="left"></div>
<div class="mid">

<hr>
<br>
<form method="POST">
Change: <input name="title"><br>
Description: <br>
<textarea cols="50" rows="20" name="description"></textarea>
<input type=submit name=submit value=add>
</form>
</div>
<div class="right"></div>
</div>
</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

// static assets

func serveTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/css")
	w.Write(theme)
}

//go:embed assets/theme.css
var theme []byte

func serveTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/javascript")
	w.Write(tools)
}

//go:embed assets/tools.js
var tools []byte
