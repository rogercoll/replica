package auth

import "github.com/rogercoll/replica"

type Creator func() replica.Authenticator

var Auths = map[string]Creator{}

func Add(name string, creator Creator) {
	Auths[name] = creator
}
