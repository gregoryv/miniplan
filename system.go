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
	tpl.Execute(w, m)
}

var index = `
<!doctype html>

<html>
<head><title>miniplan</title></head>
<body>
<ul>
{{range .Changes}}
<li>{{.Title}}</li>
{{end}}
</ul>
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
