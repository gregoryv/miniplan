package webui

import (
	"net/http"
	"text/template"

	"github.com/gregoryv/miniplan"
	. "github.com/gregoryv/miniplan"
)

func NewUI(sys *System) *UI {
	return &UI{System: sys}
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
<head><title>miniplan</title></head>
<body>
<pre>
{{range .Changes}}
{{.Title}}<b>{{.Ref}}</b>{{end}}
</pre>

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

 <style>
      body { max-width: 21cm }
  b { float: right; color: lightgray; font-weight: normal; font-size: 16px}
  .right { float: right }
  b:hover { color: blue; cursor: pointer; text-decoration: none }
  pre { line-height: 1.3em }
  .sprint { text-align: right; border-bottom: 1px dashed #000; width: 100% }
  .sprint > div { font-family: monospace; float: right; background-color: #fff; margin-top: -10px; padding-left: 30px }
  pre > div:hover { background-color: #f2f2f2 }
  strike { color: #e2e2e2 }
  strike:hover { color: #000000 }  
</style>

</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))
