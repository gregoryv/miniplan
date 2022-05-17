package miniplan

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

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
		me.InsertChange.Exec(
			uuid.Must(uuid.NewRandom()).String(),
			r.PostFormValue("title"),
			r.PostFormValue("description"),
		)
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
<input type=submit>
</form>
</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

var changesTbl = struct {
	CREATE, INSERT string
}{
	`CREATE TABLE changes (
        uuid VARCHAR(36) NULL,
        title VARCHAR(64) NULL,
        description VARCHAR(2048) NULL
    )`,
	"INSERT INTO changes(uuid, title, description) values(?,?,?)",
}

type Change struct {
	uuid.UUID
	Title
	Description
}

func (me *Change) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}

func NewPlanDB(filename string) (*PlanDB, error) {
	db, err := sql.Open("sqlite3", filename)
	if _, err = db.Exec(changesTbl.CREATE); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}

	mdb := &PlanDB{DB: db}
	stmt, err := db.Prepare(changesTbl.INSERT)
	mdb.InsertChange = stmt
	return mdb, err
}

type PlanDB struct {
	*sql.DB

	InsertChange *sql.Stmt
}

type Title string
type Description string
