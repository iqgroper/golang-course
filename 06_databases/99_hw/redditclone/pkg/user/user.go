package user

type User struct {
	ID       string `json:"id"`
	Login    string `json:"username"`
	Password string `json:"-"`
}

type NewUser struct {
	Username string
	Password string
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	Register(login, pass string) (*User, error)
}
