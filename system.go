package miniplan

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func NewDemo(dir string) (*System, func()) {
	db, _ := NewPlanDB(filepath.Join(dir, "demo.db"))
	sys := NewSystem(dir)
	sys.PlanDB = db
	sys.Create(&Change{
		Title:       "Create new changes",
		Description: "Simple todo list",
	})
	return sys, func() { db.Close(); os.RemoveAll(dir) }
}

func NewSystem(dir string) *System {
	return &System{
		rootdir:   dir,
		ViewOrder: make([]*Change, 0),
	}
}

type System struct {
	rootdir string
	*PlanDB

	ViewOrder []*Change
}

func (me *System) SetDatabase(db *PlanDB) {
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
	me.ViewOrder = changes
	rows.Close()
}

func (me *System) Save() error {
	w, err := os.Create(filepath.Join(me.rootdir, "index.json"))
	if err != nil {
		return err
	}
	defer w.Close()
	log.Print(me.ViewOrder[0])
	return json.NewEncoder(w).Encode(me.ViewOrder)
}

func (me *System) Create(v interface{}) error {
	switch v := v.(type) {
	case *Change:
		v.UUID = uuid.Must(uuid.NewRandom())
		me.ViewOrder = append(me.ViewOrder, v)
	}
	defer func() {
		if err := me.Save(); err != nil {
			log.Print(err)
		}
	}()

	return me.insert(v)
}

func (me *System) Remove(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	_, err := me.DeleteChange.Exec("%" + ref)
	return err
}

func (me *System) Update(ref string, c *Change) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	n, err := me.UpdateChange.Exec(c.Title, c.Description, "%"+ref)
	log.Println("rows affected", n)
	return err
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
