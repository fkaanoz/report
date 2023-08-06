package handlers

import (
	"encoding/json"
	"github.com/fkaanoz/middlewares/internal/app"
	"net/http"
)

// Respond will be used for successful responds. Erroneous responds will be handled by error middleware.
func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	jm, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(app.ErrorServer.Error()))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jm)
}
