package controllers

import "net/http"

func TestDebug(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test debug end point"))
}
