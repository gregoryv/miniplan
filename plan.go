package miniplan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func NewDemo(dir string) (*Plan, func()) {
	db, _ := NewPlanDB(filepath.Join(dir, "demo.db"))
	sys := NewPlan(dir)
	sys.PlanDB = db
	sys.Create(&Change{
		Title:       "Create new changes",
		Description: "Simple todo list",
	})
	return sys, func() { db.Close(); os.RemoveAll(dir) }
}

func NewPlan(dir string) *Plan {
	return &Plan{
		rootdir: dir,
		Changes: make([]*Change, 0),
	}
}

type Plan struct {
	rootdir string
	*PlanDB

	Changes []*Change
}

func (me *Plan) SetDatabase(db *PlanDB) {
	me.PlanDB = db
	// load data into memory
	rows, _ := db.Query("SELECT * FROM changes")
	var changes []*Change
	for rows.Next() {
		var c Change
		rows.Scan(&c.UUID, &c.Title, &c.Description)
		changes = append(changes, &c)
		log.Printf("load %s %s", c.Ref(), c.Title)
	}
	me.Changes = changes
	rows.Close()
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
	defer func() {
		if err := me.Save(); err != nil {
			log.Print(err)
		}
	}()

	return me.insert(v)
}

func (me *Plan) Remove(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	if _, err := me.DeleteChange.Exec("%" + ref); err != nil {
		return err
	}
	return me.Save()
}

func (me *Plan) Update(ref string, c *Change) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	_, err := me.UpdateChange.Exec(c.Title, c.Description, "%"+ref)
	if err != nil {
		return err
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
