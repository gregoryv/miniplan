package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
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

	t.Run("Delete", func(t *testing.T) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(map[string]any{
			"IdSuffix": "0091",
		})
		w, r := newRequest("DELETE", "/", &buf)

		api.Delete(w, r)
		got := w.Result()
		if err := expStatus(200, got); err != nil {
			t.Fatal(err)
		}
		if err := expJson(got); err != nil {
			t.Error(err)
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

func expStatus(exp int, got *http.Response) error {
	if got.StatusCode != exp {
		dump, _ := httputil.DumpResponse(got, true)
		return fmt.Errorf("%s\nexpected status %v", string(dump), exp)
	}
	return nil
}

func expJson(got *http.Response) error {
	var x map[string]interface{}
	err := json.NewDecoder(got.Body).Decode(&x)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
