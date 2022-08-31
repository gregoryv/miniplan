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
</head>
<body id="body">

{{range .Changes}}
<div class="entry">

<a href="#{{.Ref}}" class="idref">#</a>

<form method="POST">
<input type=hidden name="uuid" value="{{.Ref}}">
<input type=hidden name=submit value=delete>
<input type=submit value=D>
</form>


<a name="{{.Ref}}">{{.Ref}}</a>
{{.Title}}<br>


<p>{{.Description}}</p>
</div>
{{end}}


<hr>
<form method="POST">
Change: <input name="title"><br>
Description: <br>
<textarea cols="50" rows="20" name="description"></textarea>
<input type=submit name=submit value=add>
</form>

<form method="POST">
Ref: <input name="uuid"><input type=submit name=submit value=delete>
</form>

</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

func serveTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/css")
	w.Write(theme)
}

//go:embed assets/theme.css
var theme []byte
