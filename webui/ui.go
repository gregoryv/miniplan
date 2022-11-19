package webui

import (
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gregoryv/miniplan"
	. "github.com/gregoryv/miniplan"
)

func NewUI(sys *Plan) *UI {
	ui := &UI{
		Plan: sys,
	}

	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", serveAsset("assets/favicon.ico"))
	r.HandleFunc("/static/theme.css", serveAsset("assets/theme.css"))
	r.HandleFunc("/static/enhance.js", serveAsset("assets/enhance.js"))
	r.HandleFunc("/removed", ui.serveRemoved).Methods("GET")
	r.HandleFunc("/removed", audit(ui.editRemoved)).Methods("POST")
	r.HandleFunc("/", ui.serveEdit).Methods("GET")
	r.HandleFunc("/", audit(ui.editPlan)).Methods("POST")
	http.Handle("/", r)

	ui.Router = r
	return ui
}

type UI struct {
	*miniplan.Plan

	*mux.Router
}

func (me *UI) serveEdit(w http.ResponseWriter, r *http.Request) {
	var entries []EntryView
	r.ParseForm()
	filter := r.Form["filter"]
	tabber := newTabber()
	// find all tags
	tagCount := make(map[string]int)
	for _, e := range me.Entries {
		for _, tag := range e.Tags() {
			tagCount[tag]++
		}
	}
	tags := make([]string, 0, len(tagCount))
	for k, _ := range tagCount {
		tags = append(tags, k)
	}
	sort.Strings(tags)
	filters := make([]FilterView, 0, len(tags))
	for _, k := range tags {
		v := FilterView{Tag: k}
	loop:
		for _, f := range r.Form["filter"] {
			if k == f {
				v.Selected = true
				break loop
			}
		}
		filters = append(filters, v)
	}
	// list entries, using filters
	for i, e := range me.Entries {
		if len(filter) == 0 || e.HasTag(filter) {
			v := EntryView{
				Entry:   *e,
				Index:   i + 1,
				nextTab: tabber,
			}
			// calculate middle prio between previous and current
			v.InsertPrio = e.Priority + 1 // ie. above

			entries = append(entries, v)
		}
	}
	m := map[string]interface{}{
		"Entries":      entries,
		"LastPriority": 0,
		"RemovedHref":  "/removed",
		"RemovedCount": len(me.Removed),
		"Filters":      filters,
	}

	if err := edit.Execute(w, m); err != nil {
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
	// group removed entries
	groups := make([]GroupView, 3)
	groups[0].Title = "today"
	groups[1].Title = "yesterday"
	groups[2].Title = "before that"

	for i, c := range me.Removed {
		v := EntryView{
			Entry: *c,
			Index: i + 1,
		}

		switch {
		case wasToday(c.RemovedOn):
			v.When = v.RemovedOn.Format("15:04:05")
			groups[0].Entries = append(groups[0].Entries, v)

		case wasYesterday(c.RemovedOn):
			v.When = v.RemovedOn.Format("15:04")
			groups[1].Entries = append(groups[1].Entries, v)

		default:
			v.When = v.RemovedOn.Format("2006-01-02")
			groups[2].Entries = append(groups[2].Entries, v)
		}
	}
	m := map[string]interface{}{
		"RemovedGroups": groups,
		"RemovedHref":   "/removed",
		"RemovedCount":  len(me.Removed),
	}
	if err := removed.Execute(w, m); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

func wasToday(t time.Time) bool {
	if t.IsZero() {
		return false
	}
	a := time.Now().Truncate(24 * time.Hour)
	b := t.Truncate(24 * time.Hour)
	return a == b
}

func wasYesterday(t time.Time) bool {
	if t.IsZero() {
		return false
	}
	a := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -1)
	b := t.Truncate(24 * time.Hour)
	return a == b
}

func (me *UI) editRemoved(w http.ResponseWriter, r *http.Request) {
	switch r.PostFormValue("submit") {
	case "toggleDone":
		if err := me.ToggleDone(r.PostFormValue("uuid")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
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

// ----------------------------------------

type FilterView struct {
	Tag      string
	Selected bool
}

// ----------------------------------------

type GroupView struct {
	Title   string
	Entries []EntryView
}

// ----------------------------------------

type EntryView struct {
	Entry

	InsertPrio int
	Index      int
	When       string

	nextTab func() int
}

func (me *EntryView) TagClassNames() string {
	return strings.Join(me.TagNames(), " ")
}

func (me *EntryView) TagNames() []string {
	tags := me.Entry.Tags()
	names := make([]string, len(tags))
	for i := 0; i < len(tags); i++ {
		names[i] = fmt.Sprintf("tag-%s", tags[i][1:]) // strip # and prefix with tag-
	}
	return names
}

func (me *EntryView) RemovedAgo() string {
	day := 24 * time.Hour
	week := 7 * day
	age := time.Since(me.RemovedOn)
	switch {
	case me.RemovedOn.IsZero():
		return ""
	case age < time.Minute:
		return fmt.Sprintf("%vs ago", age.Truncate(time.Second).Seconds())

	case age < time.Hour:
		return fmt.Sprintf("%vm ago", age.Truncate(time.Minute).Minutes())

	case age < day:
		return fmt.Sprintf("%vh ago", age.Truncate(time.Hour).Hours())

	case age < week:
		return fmt.Sprintf("%vdays ago", int(age.Truncate(day).Hours()/24))

	default:
		return me.RemovedOn.Format("2006-01-02")
	}
}

func (me *EntryView) NextTab() int {
	return me.nextTab()
}

func (me *EntryView) LineHeight() int {
	textLen := len(me.Description)/65 + 1
	lines := strings.Count(me.Description, "\n") + 1
	if lines > textLen {
		return lines
	}
	return textLen
}

// ----------------------------------------

// newTabber returns an index counter that increases for each call
func newTabber() func() int {
	var v int
	return func() int {
		v++
		return v
	}
}

// ----------------------------------------

var (
	//go:embed edit.html
	editHtml string
	edit     = template.Must(template.New("").Parse(editHtml))
)
var (
	//go:embed removed.html
	removedHtml string
	removed     = template.Must(template.New("").Parse(removedHtml))
)

// static assets

func serveAsset(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(
			"content-type", mimeTypes[filepath.Ext(filename)],
		)
		fh, err := assets.Open(filename)
		if err != nil {
			log.Print(err)
		}
		defer fh.Close()
		io.Copy(w, fh)
	}
}

var mimeTypes = map[string]string{
	".js":  "text/javascript",
	".ico": "image/x-icon",
	".css": "text/css",
}

//go:embed assets
var assets embed.FS
