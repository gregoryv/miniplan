package webui

import (
	_ "embed"
	"net/http"
	"strings"

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
		var changes []ChangeView
		for rows.Next() {
			var c Change
			rows.Scan(&c.UUID, &c.Title, &c.Description)
			changes = append(changes, ChangeView{c})
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
		case "update":
			c := Change{
				Title:       r.PostFormValue("title"),
				Description: r.PostFormValue("description"),
			}
			err := me.Update(r.PostFormValue("uuid"), &c)
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

type ChangeView struct {
	Change
}

func (me *ChangeView) LineHeight() int {
	return strings.Count(me.Description, "\n") + 3
}
