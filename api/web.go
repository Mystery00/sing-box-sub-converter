package api

import (
	"net/http"
	"sing-box-sub-converter/server"
)

func Web(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(server.VercelHtml()))
	w.WriteHeader(http.StatusOK)
}
