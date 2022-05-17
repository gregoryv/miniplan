package miniplan

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSystem(t *testing.T) {
	sys := NewSystem()
	filename := "/tmp/test.db"
	os.RemoveAll(filename)
	db, err := NewPlanDB(filename)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.Must(uuid.NewRandom()).String()
	db.InsertChange.Exec(id, "title 1", "....")
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
	if got := buf.String(); !strings.Contains(got, "title ") {
		t.Fatal(got, "\nmissing data")
	}
}
