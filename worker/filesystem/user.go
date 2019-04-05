package filesystem

func NewUser(id string) User {
	return user{id}
}

type user struct {
	id string
}

func (u user) ID() string {
	return u.id
}
