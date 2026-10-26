// Package user ...
package user

import (
	"crypto/sha256"
	"sync"
)

type User struct {
	mu         sync.Mutex
	Name       string
	Flags      flagSet
	PasswdHash string
}

var defaultUser = User{
	Name:  "default",
	Flags: flagSet(FlagNoPass),
}

func Default() *User {
	return &defaultUser
}

func (u *User) SetPassword(pswd string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	h := sha256.New()
	h.Write([]byte(pswd))
	u.PasswdHash = string(h.Sum(nil))
	u.Flags &^= FlagNoPass.fs()
}
