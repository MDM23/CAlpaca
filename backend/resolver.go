package backend

type Resolver struct{}

func (r *Resolver) GetSelf() (*User, error) {
	return &User{ID: "abc"}, nil
}
