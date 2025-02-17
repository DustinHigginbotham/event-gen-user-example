package user

type UserService struct {
	*CreateUserCommandHandler
	*UpdateUserCommandHandler
	*UserCreatedEventHandler
	*UserUpdatedEventHandler
}

func NewUserService() *UserService {
	return &UserService{
		CreateUserCommandHandler: &CreateUserCommandHandler{},
		UpdateUserCommandHandler: &UpdateUserCommandHandler{},
		UserCreatedEventHandler:  &UserCreatedEventHandler{},
		UserUpdatedEventHandler:  &UserUpdatedEventHandler{},
	}
}
