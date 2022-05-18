package miniplan

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func NewDemo(dir string) (*System, func()) {
	db, _ := NewPlanDB(filepath.Join(dir, "demo.db"))
	sys := &System{db}
	sys.Create(&Change{
		Title:       "Create new changes",
		Description: "Simple todo list",
	})
	return sys, func() { db.Close(); os.RemoveAll(dir) }
}

func NewSystem() *System {
	return &System{}
}

type System struct {
	*PlanDB
}

func (me *System) Create(v interface{}) error {
	switch v := v.(type) {
	case *Change:
		v.UUID = uuid.Must(uuid.NewRandom())
	}
	return me.insert(v)
}

func (me *System) Remove(ref string) error {
	if ref == "" {
		return fmt.Errorf("empty ref")
	}
	_, err := me.DeleteChange.Exec("%" + ref)
	return err
}

type Change struct {
	uuid.UUID
	Title       string
	Description string
}

func (me *Change) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}
