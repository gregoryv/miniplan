package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/gregoryv/miniplan"
)

func NewRouter(api *API) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", api.Read).Methods("GET")
	r.HandleFunc("/", api.Delete).Methods("DELETE")
	return r
}

func NewAPI(sys *System) *API {
	api := &API{
		System: sys,
	}
	return api
}

type API struct {
	*System
}

func (me *API) Create(w http.ResponseWriter, r *http.Request) {
	var c Change
	json.NewDecoder(r.Body).Decode(&c)
	if err := me.System.Create(&c); err != nil {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(201)
}

func (me *API) Read(w http.ResponseWriter, r *http.Request) {
	result := map[string]any{
		"data":  "todo",
		"error": fmt.Errorf("todo"),
	}
	json.NewEncoder(w).Encode(result)
}

func (me *API) Update(w http.ResponseWriter, r *http.Request) {

}

func (me *API) Delete(w http.ResponseWriter, r *http.Request) {
	var del DeleteQuery
	json.NewDecoder(r.Body).Decode(&del)
	if err := me.System.Remove(del.IdSuffix); err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
		return
	}
}

type DeleteQuery struct {
	IdSuffix string
}
