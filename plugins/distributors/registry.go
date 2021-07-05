package distributors

import "github.com/rogercoll/replica"

type Creator func() replica.Distributor

var Distributors = map[string]Creator{}

func Add(name string, creator Creator) {
	Distributors[name] = creator
}
