package api

import (
	"net/http"
	"sing-box-sub-converter/server"
)

func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Write(server.Favicon())
	w.WriteHeader(http.StatusOK)
}
