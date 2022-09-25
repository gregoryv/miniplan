package webui

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gregoryv/miniplan"
	. "github.com/gregoryv/miniplan"
)

func NewUI(sys *Plan) *UI {
	ui := &UI{
		Plan: sys,
	}

	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", serveFavicon)
	r.HandleFunc("/static/theme.css", serveTheme)
	r.HandleFunc("/static/enhance.js", serveEnhance)
	r.HandleFunc("/removed", ui.serveRemoved).Methods("GET")
	r.HandleFunc("/removed", ui.editRemoved).Methods("POST")
	r.HandleFunc("/", ui.servePlan).Methods("GET")
	r.HandleFunc("/", ui.editPlan).Methods("POST")
	http.Handle("/", r)

	ui.Router = r
	return ui
}

type UI struct {
	*miniplan.Plan

	*mux.Router
}

func (me *UI) servePlan(w http.ResponseWriter, r *http.Request) {
	var changes []EntryView
	tabber := newTabber()
	for i, c := range me.Entries {
		v := EntryView{
			Entry:   *c,
			Index:   i + 1,
			nextTab: tabber,
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
		"RemovedHref":  "/removed",
		"RemovedCount": len(me.Removed),
	}

	if err := plan.Execute(w, m); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (me *UI) editPlan(w http.ResponseWriter, r *http.Request) {
	switch r.PostFormValue("submit") {
	case "add":
		prio, err := strconv.ParseUint(r.PostFormValue("priority"), 10, 32)
		if err != nil {
			log.Print(err)
		}
		c := Entry{
			Title:       r.PostFormValue("title"),
			Description: r.PostFormValue("description"),
			Priority:    int(prio),
		}
		_ = me.Create(&c)
		http.Redirect(w, r, "/#"+c.Ref(), 303)
		return

	case "remove":
		if err := me.Remove(r.PostFormValue("uuid")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

	case "update":
		prio, _ := strconv.ParseInt(r.PostFormValue("priority"), 10, 32)
		c := Entry{
			Title:       r.PostFormValue("title"),
			Description: r.PostFormValue("description"),
			Priority:    int(prio),
		}
		err := me.Update(r.PostFormValue("uuid"), &c)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}
	http.Redirect(w, r, "/", 303)
}

func (me *UI) serveRemoved(w http.ResponseWriter, r *http.Request) {
	var changes []EntryView
	for i, c := range me.Removed {
		v := EntryView{
			Entry: *c,
			Index: i + 1,
		}
		changes = append(changes, v)
	}
	m := map[string]interface{}{
		"Removed":      changes,
		"RemovedHref":  "/removed",
		"RemovedCount": len(me.Removed),
	}
	if err := removed.Execute(w, m); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func (me *UI) editRemoved(w http.ResponseWriter, r *http.Request) {
	switch r.PostFormValue("submit") {
	case "delete":
		if err := me.Delete(r.PostFormValue("uuid")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	case "restore":
		if err := me.Restore(r.PostFormValue("uuid")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
	}

	http.Redirect(w, r, "/removed", 303)
}

type EntryView struct {
	Entry

	InsertPrio int
	Index      int

	nextTab func() int
}

func (me *EntryView) NextTab() int {
	return me.nextTab()
}

func (me *EntryView) LineHeight() int {
	return strings.Count(me.Description, "\n") + 2
}

// ----------------------------------------

func newTabber() func() int {
	var v int
	return func() int {
		v++
		return v
	}
}

// ----------------------------------------

var (
	//go:embed plan.html
	planHtml string
	plan     = template.Must(template.New("").Parse(planHtml))
)
var (
	//go:embed removed.html
	removedHtml string
	removed     = template.Must(template.New("").Parse(removedHtml))
)

// static assets

func serveFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "image/x-icon")
	w.Write(faviconIco)
}

//go:embed assets/favicon.ico
var faviconIco []byte

func serveTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/css")
	w.Write(theme)
}

//go:embed assets/theme.css
var theme []byte

func serveEnhance(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/javascript")
	w.Write(enhance)
}

//go:embed assets/enhance.js
var enhance []byte
