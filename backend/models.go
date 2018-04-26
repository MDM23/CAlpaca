package backend

type CA struct {
}

type Organization struct {
	ID   string
	Name string
}

type User struct {
	ID        string
	Email     string
	Firstname string
	Lastname  string
}
