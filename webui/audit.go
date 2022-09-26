package webui

import (
	"log"
	"net/http"
	"time"
)

func audit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		action := r.PostFormValue("submit")
		uuid := r.PostFormValue("uuid")
		log.Printf("%s %v %s %s %v", r.RemoteAddr, r.URL.String(),
			action, uuid, time.Since(start),
		)
	}
}
