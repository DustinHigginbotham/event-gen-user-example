// Code generated by event-source-gen. DO NOT EDIT.
package app

// User
//
// A user
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserBuilder struct {
	Manager UserEventManager
}

func (a *UserBuilder) Create() *User {
	return &User{}
}

func (a *UserBuilder) Apply(entity *User, event Event) {
	switch e := event.Data.(type) {
	case *UserCreated:
		a.Manager.HandleUserCreated(entity, *e)
	case *UserUpdated:
		a.Manager.HandleUserUpdated(entity, *e)
	}
}

func (a *UserBuilder) Rebuild(events []Event) *User {
	e := a.Create()
	for _, event := range events {
		a.Apply(e, event)
	}
	return e
}
