package v1

import (
	"net/http"
)

func LoginAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Write([]byte("API V1, login"))
}