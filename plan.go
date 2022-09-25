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
	for i, _ = range me.Entries {
		if strings.HasSuffix(me.Entries[i].UUID.String(), ref) {
			break
		}
	}
	me.Entries[i].RemovedOn = time.Now()
	me.Removed = append([]*Entry{me.Entries[i]}, me.Removed...)
	me.Entries = append(me.Entries[:i], me.Entries[i+1:]...)
	return me.Save()
}

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
			break
		}
	}
	for _, e := range me.Entries {
		e.JustCreated = false
	}
	sort.Sort(ByPriority(me.Entries))
	me.fixPriority()
	return me.Save()
}

func (me *Plan) fixPriority() {
	// count those with priority > 0
	var count int
	for _, e := range me.Entries {
		if e.Priority == 0 {
			break
		}
		count++
	}
	step := 10
	switch {
	case count > 9:
		step = 5
	case count > 19:
		step = 3
	}
	v := count * step
	for _, e := range me.Entries {
		if e.Priority == 0 {
			break
		}
		e.Priority = v
		v -= step
	}
}

type Entry struct {
	UUID        uuid.UUID
	Title       string
	Description string
	Priority    int

	JustCreated bool      `json:",omitempty"`
	RemovedOn   time.Time `json:",omitempty"`
}

func (me *Entry) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}

type ByPriority []*Entry

func (b ByPriority) Len() int           { return len(b) }
func (b ByPriority) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPriority) Less(i, j int) bool { return b[i].Priority > b[j].Priority }
