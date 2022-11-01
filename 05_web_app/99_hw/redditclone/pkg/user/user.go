package user

type User struct {
	ID       uint   `json:"id,string"`
	Login    string `json:"username"`
	password string `json:"-"`
}

type NewUser struct {
	Username string
	Password string
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	Register(login, pass string) (*User, error)
}
