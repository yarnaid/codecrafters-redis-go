// Package user ...
package user

type User struct {
	Name  string
	Flags flagSet
}

var defaultUser = User{
	Name:  "default",
	Flags: flagSet(FlagNoPass),
}

func Default() *User {
	return &defaultUser
}
