package restapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/miniplan"
)

func TestAPI(t *testing.T) {
	sys, cleanup := miniplan.NewDemo(t.TempDir())
	defer cleanup()
	api := NewAPI(sys)

	t.Run("Create", func(t *testing.T) {
		w, r := newRequest("GET", "/", nil)
		api.Create(w, r)
		got := w.Result()
		if exp := 201; got.StatusCode != exp {
			t.Errorf("%s, expected %v", got.Status, exp)
		}
	})
}

func newRequest(method, path string, body io.Reader) (
	*httptest.ResponseRecorder, *http.Request,
) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, body)
	return w, r
}
