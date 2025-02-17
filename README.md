# Example: User from event-gen

Still a work in progress, but everything present in the `gen` directory is automatically generated from [`event-gen`](https://github.com/DustinHigginbotham/event-gen).

Exploring `cmd/users/main.go` will show you how to use the resulting system.

Additionally, `event-gen` also generates the files: `domains/user/command_create.go`, `domains/user/command_update.go`, and `domains/user/handlers.go`. The system tries its best to figure out what it can automatically generate in those files to make wiring this as easy as possible. These files are fine to change to your liking. If you run the generator again, it won't override these functions if already present.

## Running the generator

If you have go 1.24+ then you can re-run the generator like this:

```bash
go tool event-gen
```
