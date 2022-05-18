package spec

import (
	"fmt"
	"testing"
)

// ----------------------------------------

type R struct {
	ref    string
	Short  string
	Reason string
	Rel    []*R
}

func (me *R) String() string {
	return fmt.Sprintf("%s %s", me.Short, me.ref)
}

func Reflog(t testing.TB, r *R) {
	t.Helper()
	t.Log(r.String())
}
