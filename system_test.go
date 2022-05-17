package miniplan

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSystem(t *testing.T) {
	sys := NewSystem()
	db, _ := NewPlanDB("/tmp/test.db")
	sys.PlanDB = db

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	sys.ServeHTTP(w, r)
	defer sys.Close()

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatal(resp.Status)
	}
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	if got := buf.String(); !strings.Contains(got, "<html>") {
		t.Fatal(got, "\nmissing data")
	}
}
