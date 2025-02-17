package user

import (
	"context"

	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
)

type UpdateUserCommandHandler struct{}

func (h *UpdateUserCommandHandler) UpdateUser(
	ctx context.Context,
	a *app.App,
	cmd *app.UpdateUser,
) (*app.UserUpdated, error) {
	return &app.UserUpdated{
		ID:        cmd.ID,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
	}, nil
}
