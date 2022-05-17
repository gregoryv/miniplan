package miniplan

import (
	"os"
	"testing"
)

func TestSystem(t *testing.T) {
	sys := NewSystem()
	filename := "/tmp/test.db"
	os.RemoveAll(filename)
	db, err := NewPlanDB(filename)
	if err != nil {
		t.Fatal(err)
	}
	sys.PlanDB = db

	if err := sys.Create(&Change{Title: "test"}); err != nil {
		t.Fatal(err)
	}
}
