package model

type User struct {
	ID    int
	Name  string
	Email string
}

type Error struct {
	Message string
}

func (e *Error) Error() string {
	return "SOME ERROR MESSAGE"
}
