package main

import (
	"github.com/DustinHigginbotham/event-gen-user-example/domains/user"
	app "github.com/DustinHigginbotham/event-gen-user-example/gen"
	"github.com/DustinHigginbotham/event-gen-user-example/internal/events/inmem"
	"github.com/DustinHigginbotham/event-gen-user-example/internal/events/sqlite"
)

func main() {
	sys := app.NewApp()

	queue := inmem.NewQueue(sys, 100)
	store := sqlite.New(sys)
	if err := store.Connect("file:users.s3db?cache=shared"); err != nil {
		panic(err)
	}

	defer store.Close()

	sys.SetQueue(queue)
	sys.SetStore(store)

	sys.RegisterUserService(user.NewUserService())

}
