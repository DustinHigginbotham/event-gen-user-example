package user

import (
	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
)

type UserCreatedEventHandler struct{}

func (h *UserCreatedEventHandler) HandleUserCreated(
	a *app.User,
	ev app.UserCreated,
) error {
	a.ID = ev.ID
	a.FirstName = ev.FirstName
	a.LastName = ev.LastName
	a.Email = ev.Email
	return nil
}

type UserUpdatedEventHandler struct{}

func (h *UserUpdatedEventHandler) HandleUserUpdated(
	a *app.User,
	ev app.UserUpdated,
) error {
	a.ID = ev.ID
	a.FirstName = ev.FirstName
	a.LastName = ev.LastName
	return nil
}