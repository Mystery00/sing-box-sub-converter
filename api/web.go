package api

import (
	"net/http"
)

func Web(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexHtml))
}
