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
	"instance.version":          {Name: "instance.version", Method: "GET", Path: "/api/instance/version", Mutation: ReadOnly},
	"instance.status":           {Name: "instance.status", Method: "GET", Path: "/api/instance", Mutation: ReadOnly},
	"accounts.list":             {Name: "accounts.list", Method: "GET", Path: "/api/accounts", Mutation: ReadOnly},
	"dns.nameservers.list":      {Name: "dns.nameservers.list", Method: "GET", Path: "/api/dns/nameservers", Mutation: ReadOnly},
	"dns.nameservers.get":       {Name: "dns.nameservers.get", Method: "GET", Path: "/api/dns/nameservers/{nsgroupId}", Mutation: ReadOnly},
	"dns.settings":              {Name: "dns.settings", Method: "GET", Path: "/api/dns/settings", Mutation: ReadOnly},
	"dns.zones.list":            {Name: "dns.zones.list", Method: "GET", Path: "/api/dns/zones", Mutation: ReadOnly},
	"dns.zones.get":             {Name: "dns.zones.get", Method: "GET", Path: "/api/dns/zones/{zoneId}", Mutation: ReadOnly},
	"identity_providers.list":   {Name: "identity_providers.list", Method: "GET", Path: "/api/identity-providers", Mutation: ReadOnly},
	"identity_providers.get":    {Name: "identity_providers.get", Method: "GET", Path: "/api/identity-providers/{idpId}", Mutation: ReadOnly},
	"posture_checks.list":       {Name: "posture_checks.list", Method: "GET", Path: "/api/posture-checks", Mutation: ReadOnly},
	"posture_checks.get":        {Name: "posture_checks.get", Method: "GET", Path: "/api/posture-checks/{postureCheckId}", Mutation: ReadOnly},
	"events.audit":              {Name: "events.audit", Method: "GET", Path: "/api/events/audit", Mutation: ReadOnly},
	"setup_keys.list":           {Name: "setup_keys.list", Method: "GET", Path: "/api/setup-keys", Mutation: ReadOnly},
	"setup_keys.get":            {Name: "setup_keys.get", Method: "GET", Path: "/api/setup-keys/{keyId}", Mutation: ReadOnly},
	"locations.countries":       {Name: "locations.countries", Method: "GET", Path: "/api/locations/countries", Mutation: ReadOnly},
	"locations.cities":          {Name: "locations.cities", Method: "GET", Path: "/api/locations/countries/{country}/cities", Mutation: ReadOnly},
	"users.tokens.list":         {Name: "users.tokens.list", Method: "GET", Path: "/api/users/{userId}/tokens", Mutation: ReadOnly},
	"users.tokens.get":          {Name: "users.tokens.get", Method: "GET", Path: "/api/users/{userId}/tokens/{tokenId}", Mutation: ReadOnly},
	"networks.resources.list":   {Name: "networks.resources.list", Method: "GET", Path: "/api/networks/{networkId}/resources", Mutation: ReadOnly},
	"networks.resources.get":    {Name: "networks.resources.get", Method: "GET", Path: "/api/networks/{networkId}/resources/{resourceId}", Mutation: ReadOnly},
	"networks.routers.list":     {Name: "networks.routers.list", Method: "GET", Path: "/api/networks/{networkId}/routers", Mutation: ReadOnly},
	"networks.routers.get":      {Name: "networks.routers.get", Method: "GET", Path: "/api/networks/{networkId}/routers/{routerId}", Mutation: ReadOnly},
	"networks.routers.list_all": {Name: "networks.routers.list_all", Method: "GET", Path: "/api/networks/routers", Mutation: ReadOnly},
	"networks.list":             {Name: "networks.list", Method: "GET", Path: "/api/networks", Mutation: ReadOnly},
	"networks.get":              {Name: "networks.get", Method: "GET", Path: "/api/networks/{networkId}", Mutation: ReadOnly},
	"routes.list":               {Name: "routes.list", Method: "GET", Path: "/api/routes", Mutation: ReadOnly},
	"routes.get":                {Name: "routes.get", Method: "GET", Path: "/api/routes/{routeId}", Mutation: ReadOnly},
	"users.current":             {Name: "users.current", Method: "GET", Path: "/api/users/current", Mutation: ReadOnly},
	"users.list":                {Name: "users.list", Method: "GET", Path: "/api/users", Mutation: ReadOnly},
	"users.invites":             {Name: "users.invites", Method: "GET", Path: "/api/users/invites", Mutation: ReadOnly},
	"groups.list":               {Name: "groups.list", Method: "GET", Path: "/api/groups", Mutation: ReadOnly},
	"groups.get":                {Name: "groups.get", Method: "GET", Path: "/api/groups/{groupId}", Mutation: ReadOnly},
	"groups.update":             {Name: "groups.update", Method: "PUT", Path: "/api/groups/{groupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.list":                {Name: "peers.list", Method: "GET", Path: "/api/peers", Mutation: ReadOnly},
	"peers.get":                 {Name: "peers.get", Method: "GET", Path: "/api/peers/{peerId}", Mutation: ReadOnly},
	"policies.list":             {Name: "policies.list", Method: "GET", Path: "/api/policies", Mutation: ReadOnly},
	"policies.get":              {Name: "policies.get", Method: "GET", Path: "/api/policies/{policyId}", Mutation: ReadOnly},
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
