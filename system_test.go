package miniplan

import (
	"testing"
)

func TestSystem(t *testing.T) {
	sys, cleanup := NewDemo(t.TempDir())
	defer cleanup()

	if err := sys.Create(&Change{Title: "test"}); err != nil {
		t.Fatal(err)
	}
}
