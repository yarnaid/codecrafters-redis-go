// Package user ...
package user

import (
	"crypto/sha256"
	"fmt"
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

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (u *User) SetPassword(pswd string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.PasswdHash = hash(pswd)
	u.Flags &^= FlagNoPass.fs()
}

func (u *User) CheckPassword(pswd string) bool {
	h := hash(pswd)
	return h == u.PasswdHash
}
