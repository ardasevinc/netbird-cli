package operations

import (
	"fmt"
	"sort"
)

type MutationClass string

const (
	ReadOnly      MutationClass = "read_only"
	Consequential MutationClass = "consequential"
	Cosmetic      MutationClass = "cosmetic"
)

type Definition struct {
	Name               string
	Method             string
	Path               string
	Mutation           MutationClass
	RequiresPreimage   bool
	RequiresReadBack   bool
	DispatcherAdmitted bool
}

var definitions = map[string]Definition{
	"instance.version": {Name: "instance.version", Method: "GET", Path: "/api/instance/version", Mutation: ReadOnly},
	"instance.status":  {Name: "instance.status", Method: "GET", Path: "/api/instance", Mutation: ReadOnly},
	"accounts.list":    {Name: "accounts.list", Method: "GET", Path: "/api/accounts", Mutation: ReadOnly},
	"users.current":    {Name: "users.current", Method: "GET", Path: "/api/users/current", Mutation: ReadOnly},
	"users.list":       {Name: "users.list", Method: "GET", Path: "/api/users", Mutation: ReadOnly},
	"users.invites":    {Name: "users.invites", Method: "GET", Path: "/api/users/invites", Mutation: ReadOnly},
	"groups.list":      {Name: "groups.list", Method: "GET", Path: "/api/groups", Mutation: ReadOnly},
	"groups.get":       {Name: "groups.get", Method: "GET", Path: "/api/groups/{groupId}", Mutation: ReadOnly},
	"groups.update":    {Name: "groups.update", Method: "PUT", Path: "/api/groups/{groupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.list":       {Name: "peers.list", Method: "GET", Path: "/api/peers", Mutation: ReadOnly},
	"peers.get":        {Name: "peers.get", Method: "GET", Path: "/api/peers/{peerId}", Mutation: ReadOnly},
	"policies.list":    {Name: "policies.list", Method: "GET", Path: "/api/policies", Mutation: ReadOnly},
	"policies.get":     {Name: "policies.get", Method: "GET", Path: "/api/policies/{policyId}", Mutation: ReadOnly},
}

func Lookup(name string) (Definition, error) {
	definition, ok := definitions[name]
	if !ok {
		return Definition{}, fmt.Errorf("operation %q is not admitted in the coverage registry", name)
	}
	return definition, nil
}

func All() []Definition {
	result := make([]Definition, 0, len(definitions))
	for _, definition := range definitions {
		result = append(result, definition)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
