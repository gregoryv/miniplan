package webui

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gregoryv/miniplan"
	. "github.com/gregoryv/miniplan"
)

func NewUI(sys *Plan) *UI {
	http.HandleFunc("/static/theme.css", serveTheme)
	http.HandleFunc("/static/tools.js", serveTools)

	ui := &UI{Plan: sys}
	http.Handle("/", ui)
	return ui
}

type UI struct {
	*miniplan.Plan
}

func (me *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var changes []ChangeView
		for i, c := range me.Entries {
			v := ChangeView{
				Entry: *c,
			}
			// calculate middle prio between previous and current
			v.InsertPrio = c.Priority + 10 // ie. above
			if i > 0 {
				diff := (me.Entries[i-1].Priority - c.Priority) / 2
				v.InsertPrio = c.Priority + diff
			}
			changes = append(changes, v)
		}
		m := map[string]interface{}{
			"Changes":      changes,
			"LastPriority": 0,
		}
		err := tpl.Execute(w, m)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

	case "POST":
		switch r.PostFormValue("submit") {
		case "add":
			prio, err := strconv.ParseUint(r.PostFormValue("priority"), 10, 32)
			if err != nil {
				log.Print(err)
			}
			c := Entry{
				Title:       r.PostFormValue("title"),
				Description: r.PostFormValue("description"),
				Priority:    uint32(prio),
			}
			_ = me.Create(&c)

		case "delete":
			err := me.Remove(r.PostFormValue("uuid"))
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
			}
		case "update":
			prio, _ := strconv.ParseUint(r.PostFormValue("priority"), 10, 32)
			c := Entry{
				Title:       r.PostFormValue("title"),
				Description: r.PostFormValue("description"),
				Priority:    uint32(prio),
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
	Entry

	InsertPrio uint32
}

func (me *ChangeView) LineHeight() int {
	return strings.Count(me.Description, "\n") + 2
}

// ----------------------------------------

var tpl = template.Must(template.New("").Parse(indexHtml))

//go:embed index.html
var indexHtml string

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
