package miniplan

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func NewSystemLite(filename string) (*System, error) {
	db, err := sql.Open("sqlite3", filename)
	sys := NewSystem()
	sys.DB = db
	return sys, err
}

func NewSystem() *System {
	return &System{}
}

type System struct {
	*sql.DB
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
