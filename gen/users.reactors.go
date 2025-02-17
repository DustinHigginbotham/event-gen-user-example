package app

import (
	"context"
)

type UserReactorManager interface {
	HandleWelcomeEmailReactor(ctx context.Context, app *App, event *UserCreated)
}
