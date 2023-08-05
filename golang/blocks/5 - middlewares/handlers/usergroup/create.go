package usergroup

import (
	"context"
	"encoding/json"
	"github.com/fkaanoz/middlewares/handlers"
	"github.com/fkaanoz/middlewares/models"
	"io"
	"net/http"
)

// Create returns an error. If there is something bad, error middleware will handle it and respond to the client.
func Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var u models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &u)
	if err != nil {
		return err
	}

	users = append(users, u)

	handlers.Respond(w, http.StatusCreated, u)

	return nil
}
