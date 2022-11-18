package miniplan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewDemo(dir string) (*Plan, func()) {
	p := NewPlan(filepath.Join(dir, "demo.json"))
	p.Create(&Entry{
		Title:       "Create new changes",
		Description: "Simple todo list",
	})
	return p, func() { os.RemoveAll(dir) }
}

func NewPlan(planfile string) *Plan {
	return &Plan{
		planfile: planfile,
		Entries:  make([]*Entry, 0),
		Removed:  make([]*Entry, 0),
	}
}

type Plan struct {
	planfile string

	Entries []*Entry
	Removed []*Entry
}

func (me *Plan) Load() {
	// load data into memory
	fh, err := os.Open(me.planfile)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	if err := json.NewDecoder(fh).Decode(me); err != nil {
		log.Print("Load ", err)
	}
}

func (me *Plan) Save() error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(me); err != nil {
		return err
	}

	var tidy bytes.Buffer
	json.Indent(&tidy, buf.Bytes(), "", "  ")
	w, err := os.Create(me.planfile)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, &tidy)
	return err
}

func (me *Plan) Create(v *Entry) error {
	v.UUID = uuid.Must(uuid.NewRandom())
	me.Entries = append(me.Entries, v)
	sort.Sort(ByPriority(me.Entries))
	for _, e := range me.Entries {
		e.JustCreated = false
	}
	v.JustCreated = true
	return me.Save()
}

func (me *Plan) Remove(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	var i int
	var e *Entry
	for i, e = range me.Entries {
		if strings.HasSuffix(me.Entries[i].UUID.String(), ref) {
			break
		}
	}
	e.RemovedOn = time.Now()
	me.Removed = append([]*Entry{e}, me.Removed...)
	me.Entries = append(me.Entries[:i], me.Entries[i+1:]...)
	return me.Save()
}

func (me *Plan) ToggleDone(ref string) error {
	e, err := me.findEntry(me.Removed, ref)
	if err != nil {
		return err
	}
	e.Done = !e.Done
	return me.Save()
}

func (me *Plan) findEntry(list []*Entry, ref string) (*Entry, error) {
	if ref == "" {
		return nil, fmt.Errorf("empty ref")
	}
	for _, e := range me.Removed {
		if strings.HasSuffix(e.UUID.String(), ref) {
			return e, nil
		}
	}
	return nil, ErrNotFound
}

var ErrNotFound = fmt.Errorf("not found")

func (me *Plan) Delete(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	var i int
	for i, _ = range me.Removed {
		if strings.HasSuffix(me.Removed[i].UUID.String(), ref) {
			break
		}
	}
	me.Removed = append(me.Removed[:i], me.Removed[i+1:]...)
	return me.Save()
}

func (me *Plan) Restore(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	var i int
	for i, _ = range me.Removed {
		if strings.HasSuffix(me.Removed[i].UUID.String(), ref) {
			break
		}
	}
	entry := me.Removed[i]
	entry.Priority = 0 // put it last
	entry.RemovedOn = time.Time{}

	me.Entries = append(me.Entries, entry)
	me.Removed = append(me.Removed[:i], me.Removed[i+1:]...) // clear
	return me.Save()
}

func (me *Plan) Update(ref string, in *Entry) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	for _, c := range me.Entries {
		if strings.HasSuffix(c.UUID.String(), ref) {
			c.Title = in.Title
			c.Description = in.Description
			c.Priority = in.Priority
		} else {
			// push down other entry
			if c.Priority == in.Priority {
				c.Priority--
			}
		}
	}
	// reset any just created entry
	for _, e := range me.Entries {
		e.JustCreated = false
	}
	sort.Sort(ByPriority(me.Entries))
	me.fixPriority()
	return me.Save()
}

func (me *Plan) fixPriority() {
	v := len(me.Entries)
	for _, e := range me.Entries {
		e.Priority = v
		v -= 1
	}
}
