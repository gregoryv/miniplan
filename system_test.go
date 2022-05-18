package miniplan

import (
	"testing"

	. "github.com/gregoryv/miniplan/spec"
)

func TestSystem(t *testing.T) {
	sys, cleanup := NewDemo(t.TempDir())
	defer cleanup()

	c := Change{
		Title: "Something new...",
	}
	t.Run("Create", func(t *testing.T) {
		if err := sys.Create(&c); err != nil {
			t.Fatal(err)
		}
		if c.UUID.ID() == 0 {
			t.Fatal("uuid not set")
		}
	})

	t.Run("Remove", func(t *testing.T) {
		Reflog(t, YRA05)
		if err := sys.Remove(c.Ref()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Remove empty", func(t *testing.T) {
		if err := sys.Remove(""); err == nil {
			t.Error("should fail")
		}
	})
}
