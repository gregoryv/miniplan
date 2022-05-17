package miniplan

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func NewSystem() *System {
	return &System{}
}

type System struct {
	*PlanDB
}

func (me *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, me)
}

var index = `
<!doctype html>

<html>
<body>

</body>
</html>
`

var tpl = template.Must(template.New("").Parse(index))

var changesTbl = struct {
	CREATE, INSERT string
}{
	`CREATE TABLE changes (
        title VARCHAR(64) NULL,
        description VARCHAR(2048) NULL
    )`,
	"INSERT INTO changes(title, description) values(?,?)",
}

type Change struct {
	Title       string
	Description string
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
