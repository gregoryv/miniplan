package miniplan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func NewDemo(dir string) (*Plan, func()) {
	sys := NewPlan(dir)
	sys.Create(&Change{
		Title:       "Create new changes",
		Description: "Simple todo list",
	})
	return sys, func() { os.RemoveAll(dir) }
}

func NewPlan(dir string) *Plan {
	p := &Plan{
		rootdir: dir,
		Changes: make([]*Change, 0),
	}
	return p
}

type Plan struct {
	rootdir string
	*PlanDB

	Changes []*Change
}

func (me *Plan) Load() {
	// load data into memory
	fh, err := os.Open(filepath.Join(me.rootdir, "index.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	if err := json.NewDecoder(fh).Decode(&me.Changes); err != nil {
		log.Print("Load ", err)
	}
}

func (me *Plan) Save() error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(me.Changes); err != nil {
		return err
	}

	var tidy bytes.Buffer
	json.Indent(&tidy, buf.Bytes(), "", "  ")
	w, err := os.Create(filepath.Join(me.rootdir, "index.json"))
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, &tidy)
	return err
}

func (me *Plan) Create(v interface{}) error {
	switch v := v.(type) {
	case *Change:
		v.UUID = uuid.Must(uuid.NewRandom())
		me.Changes = append(me.Changes, v)
	}
	return me.Save()
}

func (me *Plan) Remove(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	var i int
	for i, _ = range me.Changes {
		if strings.HasSuffix(me.Changes[i].UUID.String(), ref) {
			break
		}
	}
	me.Changes = append(me.Changes[:i], me.Changes[i+1:]...)
	return me.Save()
}

func (me *Plan) Update(ref string, in *Change) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	for _, c := range me.Changes {
		if strings.HasSuffix(c.UUID.String(), ref) {
			c.Title = in.Title
			c.Description = in.Description
			break
		}
	}
	return me.Save()
}

type Change struct {
	UUID        uuid.UUID
	Title       string
	Description string
}

func (me *Change) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}
