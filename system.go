package miniplan

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func NewSystem() *System {
	return &System{}
}

type System struct {
	*PlanDB
}

func (me *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
				UUID:        uuid.Must(uuid.NewRandom()),
				Title:       r.PostFormValue("title"),
				Description: r.PostFormValue("description"),
			}
			_ = me.Create(&c)

		case "delete":
			_, err := me.DeleteChange.Exec(
				"%" + r.PostFormValue("uuid"),
			)
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

func (me *System) Create(v interface{}) error {
	return me.insert(v)
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
</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

type Change struct {
	uuid.UUID
	Title       string
	Description string
}

func (me *Change) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}
