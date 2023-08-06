package usergroup

import (
	"context"
	"github.com/fkaanoz/middlewares/internal/handlers"
	"github.com/fkaanoz/middlewares/internal/models"
	"net/http"
)

var users = []models.User{
	{
		ID:   1,
		Name: "fkz",
	},
}

// List returns an error. If there is something bad, error middleware will handle it and respond to the client.
func List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	handlers.Respond(w, http.StatusOK, users)
	return nil
}
