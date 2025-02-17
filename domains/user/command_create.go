package user

import (
	"context"

	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
)

type CreateUserCommandHandler struct{}

func (h *CreateUserCommandHandler) CreateUser(
	ctx context.Context,
	a *app.App,
	cmd *app.CreateUser,
) (*app.UserCreated, error) {
	return &app.UserCreated{
		ID:        cmd.ID,
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Email:     cmd.Email,
	}, nil
}
