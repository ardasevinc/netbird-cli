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
	"instance.version":                  {Name: "instance.version", Method: "GET", Path: "/api/instance/version", Mutation: ReadOnly},
	"instance.status":                   {Name: "instance.status", Method: "GET", Path: "/api/instance", Mutation: ReadOnly},
	"accounts.list":                     {Name: "accounts.list", Method: "GET", Path: "/api/accounts", Mutation: ReadOnly},
	"accounts.update":                   {Name: "accounts.update", Method: "PUT", Path: "/api/accounts/{accountId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"accounts.delete":                   {Name: "accounts.delete", Method: "DELETE", Path: "/api/accounts/{accountId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.nameservers.list":              {Name: "dns.nameservers.list", Method: "GET", Path: "/api/dns/nameservers", Mutation: ReadOnly},
	"dns.nameservers.get":               {Name: "dns.nameservers.get", Method: "GET", Path: "/api/dns/nameservers/{nsgroupId}", Mutation: ReadOnly},
	"dns.nameservers.create":            {Name: "dns.nameservers.create", Method: "POST", Path: "/api/dns/nameservers", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.nameservers.update":            {Name: "dns.nameservers.update", Method: "PUT", Path: "/api/dns/nameservers/{nsgroupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.nameservers.delete":            {Name: "dns.nameservers.delete", Method: "DELETE", Path: "/api/dns/nameservers/{nsgroupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.settings":                      {Name: "dns.settings", Method: "GET", Path: "/api/dns/settings", Mutation: ReadOnly},
	"dns.settings.update":               {Name: "dns.settings.update", Method: "PUT", Path: "/api/dns/settings", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.zones.list":                    {Name: "dns.zones.list", Method: "GET", Path: "/api/dns/zones", Mutation: ReadOnly},
	"dns.zones.get":                     {Name: "dns.zones.get", Method: "GET", Path: "/api/dns/zones/{zoneId}", Mutation: ReadOnly},
	"dns.zones.create":                  {Name: "dns.zones.create", Method: "POST", Path: "/api/dns/zones", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.zones.delete":                  {Name: "dns.zones.delete", Method: "DELETE", Path: "/api/dns/zones/{zoneId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.zones.update":                  {Name: "dns.zones.update", Method: "PUT", Path: "/api/dns/zones/{zoneId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"identity_providers.list":           {Name: "identity_providers.list", Method: "GET", Path: "/api/identity-providers", Mutation: ReadOnly},
	"identity_providers.get":            {Name: "identity_providers.get", Method: "GET", Path: "/api/identity-providers/{idpId}", Mutation: ReadOnly},
	"identity_providers.create":         {Name: "identity_providers.create", Method: "POST", Path: "/api/identity-providers", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"identity_providers.update":         {Name: "identity_providers.update", Method: "PUT", Path: "/api/identity-providers/{idpId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"identity_providers.delete":         {Name: "identity_providers.delete", Method: "DELETE", Path: "/api/identity-providers/{idpId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"reverse_proxy_tokens.create":       {Name: "reverse_proxy_tokens.create", Method: "POST", Path: "/api/reverse-proxies/proxy-tokens", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"reverse_proxy_tokens.delete":       {Name: "reverse_proxy_tokens.delete", Method: "DELETE", Path: "/api/reverse-proxies/proxy-tokens/{tokenId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"reverse_proxy_domains.create":      {Name: "reverse_proxy_domains.create", Method: "POST", Path: "/api/reverse-proxies/domains", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"reverse_proxy_domains.delete":      {Name: "reverse_proxy_domains.delete", Method: "DELETE", Path: "/api/reverse-proxies/domains/{domainId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"posture_checks.list":               {Name: "posture_checks.list", Method: "GET", Path: "/api/posture-checks", Mutation: ReadOnly},
	"posture_checks.get":                {Name: "posture_checks.get", Method: "GET", Path: "/api/posture-checks/{postureCheckId}", Mutation: ReadOnly},
	"posture_checks.create":             {Name: "posture_checks.create", Method: "POST", Path: "/api/posture-checks", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"posture_checks.update":             {Name: "posture_checks.update", Method: "PUT", Path: "/api/posture-checks/{postureCheckId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"posture_checks.delete":             {Name: "posture_checks.delete", Method: "DELETE", Path: "/api/posture-checks/{postureCheckId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"events.audit":                      {Name: "events.audit", Method: "GET", Path: "/api/events/audit", Mutation: ReadOnly},
	"events.network_traffic":            {Name: "events.network_traffic", Method: "GET", Path: "/api/events/network-traffic", Mutation: ReadOnly},
	"events.proxy":                      {Name: "events.proxy", Method: "GET", Path: "/api/events/proxy", Mutation: ReadOnly},
	"event_streaming.create":            {Name: "event_streaming.create", Method: "POST", Path: "/api/event-streaming", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"event_streaming.delete":            {Name: "event_streaming.delete", Method: "DELETE", Path: "/api/event-streaming/{id}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"event_streaming.update":            {Name: "event_streaming.update", Method: "PUT", Path: "/api/event-streaming/{id}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"setup_keys.list":                   {Name: "setup_keys.list", Method: "GET", Path: "/api/setup-keys", Mutation: ReadOnly},
	"setup_keys.get":                    {Name: "setup_keys.get", Method: "GET", Path: "/api/setup-keys/{keyId}", Mutation: ReadOnly},
	"setup_keys.create":                 {Name: "setup_keys.create", Method: "POST", Path: "/api/setup-keys", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"setup_keys.update":                 {Name: "setup_keys.update", Method: "PUT", Path: "/api/setup-keys/{keyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"setup_keys.delete":                 {Name: "setup_keys.delete", Method: "DELETE", Path: "/api/setup-keys/{keyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"locations.countries":               {Name: "locations.countries", Method: "GET", Path: "/api/locations/countries", Mutation: ReadOnly},
	"locations.cities":                  {Name: "locations.cities", Method: "GET", Path: "/api/locations/countries/{country}/cities", Mutation: ReadOnly},
	"users.tokens.list":                 {Name: "users.tokens.list", Method: "GET", Path: "/api/users/{userId}/tokens", Mutation: ReadOnly},
	"users.tokens.get":                  {Name: "users.tokens.get", Method: "GET", Path: "/api/users/{userId}/tokens/{tokenId}", Mutation: ReadOnly},
	"users.tokens.delete":               {Name: "users.tokens.delete", Method: "DELETE", Path: "/api/users/{userId}/tokens/{tokenId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.tokens.create":               {Name: "users.tokens.create", Method: "POST", Path: "/api/users/{userId}/tokens", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.create":                      {Name: "users.create", Method: "POST", Path: "/api/users", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.update":                      {Name: "users.update", Method: "PUT", Path: "/api/users/{userId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.delete":                      {Name: "users.delete", Method: "DELETE", Path: "/api/users/{userId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.approve":                     {Name: "users.approve", Method: "POST", Path: "/api/users/{userId}/approve", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.reject":                      {Name: "users.reject", Method: "DELETE", Path: "/api/users/{userId}/reject", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.password.update":             {Name: "users.password.update", Method: "PUT", Path: "/api/users/{userId}/password", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.invite.resend":               {Name: "users.invite.resend", Method: "POST", Path: "/api/users/{userId}/invite", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.resources.list":           {Name: "networks.resources.list", Method: "GET", Path: "/api/networks/{networkId}/resources", Mutation: ReadOnly},
	"networks.resources.get":            {Name: "networks.resources.get", Method: "GET", Path: "/api/networks/{networkId}/resources/{resourceId}", Mutation: ReadOnly},
	"networks.resources.create":         {Name: "networks.resources.create", Method: "POST", Path: "/api/networks/{networkId}/resources", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.resources.update":         {Name: "networks.resources.update", Method: "PUT", Path: "/api/networks/{networkId}/resources/{resourceId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.resources.delete":         {Name: "networks.resources.delete", Method: "DELETE", Path: "/api/networks/{networkId}/resources/{resourceId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.routers.list":             {Name: "networks.routers.list", Method: "GET", Path: "/api/networks/{networkId}/routers", Mutation: ReadOnly},
	"networks.routers.get":              {Name: "networks.routers.get", Method: "GET", Path: "/api/networks/{networkId}/routers/{routerId}", Mutation: ReadOnly},
	"networks.routers.create":           {Name: "networks.routers.create", Method: "POST", Path: "/api/networks/{networkId}/routers", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.routers.update":           {Name: "networks.routers.update", Method: "PUT", Path: "/api/networks/{networkId}/routers/{routerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.routers.delete":           {Name: "networks.routers.delete", Method: "DELETE", Path: "/api/networks/{networkId}/routers/{routerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.routers.list_all":         {Name: "networks.routers.list_all", Method: "GET", Path: "/api/networks/routers", Mutation: ReadOnly},
	"dns.records.list":                  {Name: "dns.records.list", Method: "GET", Path: "/api/dns/zones/{zoneId}/records", Mutation: ReadOnly},
	"dns.records.get":                   {Name: "dns.records.get", Method: "GET", Path: "/api/dns/zones/{zoneId}/records/{recordId}", Mutation: ReadOnly},
	"dns.records.create":                {Name: "dns.records.create", Method: "POST", Path: "/api/dns/zones/{zoneId}/records", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.records.update":                {Name: "dns.records.update", Method: "PUT", Path: "/api/dns/zones/{zoneId}/records/{recordId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"dns.records.delete":                {Name: "dns.records.delete", Method: "DELETE", Path: "/api/dns/zones/{zoneId}/records/{recordId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"ingress.peers.list":                {Name: "ingress.peers.list", Method: "GET", Path: "/api/ingress/peers", Mutation: ReadOnly},
	"ingress.peers.get":                 {Name: "ingress.peers.get", Method: "GET", Path: "/api/ingress/peers/{ingressPeerId}", Mutation: ReadOnly},
	"ingress.peers.create":              {Name: "ingress.peers.create", Method: "POST", Path: "/api/ingress/peers", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"ingress.peers.update":              {Name: "ingress.peers.update", Method: "PUT", Path: "/api/ingress/peers/{ingressPeerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"ingress.peers.delete":              {Name: "ingress.peers.delete", Method: "DELETE", Path: "/api/ingress/peers/{ingressPeerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.settings.update":     {Name: "agent_network.settings.update", Method: "PUT", Path: "/api/agent-network/settings", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.settings.create":     {Name: "agent_network.settings.create", Method: "POST", Path: "/api/agent-network/settings", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.settings.delete":     {Name: "agent_network.settings.delete", Method: "DELETE", Path: "/api/agent-network/settings", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.budget_rules.create": {Name: "agent_network.budget_rules.create", Method: "POST", Path: "/api/agent-network/budget-rules", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.budget_rules.update": {Name: "agent_network.budget_rules.update", Method: "PUT", Path: "/api/agent-network/budget-rules/{ruleId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.budget_rules.delete": {Name: "agent_network.budget_rules.delete", Method: "DELETE", Path: "/api/agent-network/budget-rules/{ruleId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.guardrails.create":   {Name: "agent_network.guardrails.create", Method: "POST", Path: "/api/agent-network/guardrails", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.guardrails.update":   {Name: "agent_network.guardrails.update", Method: "PUT", Path: "/api/agent-network/guardrails/{guardrailId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.guardrails.delete":   {Name: "agent_network.guardrails.delete", Method: "DELETE", Path: "/api/agent-network/guardrails/{guardrailId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.policies.create":     {Name: "agent_network.policies.create", Method: "POST", Path: "/api/agent-network/policies", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.policies.update":     {Name: "agent_network.policies.update", Method: "PUT", Path: "/api/agent-network/policies/{policyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.policies.delete":     {Name: "agent_network.policies.delete", Method: "DELETE", Path: "/api/agent-network/policies/{policyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.providers.create":    {Name: "agent_network.providers.create", Method: "POST", Path: "/api/agent-network/providers", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.providers.update":    {Name: "agent_network.providers.update", Method: "PUT", Path: "/api/agent-network/providers/{providerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"agent_network.providers.delete":    {Name: "agent_network.providers.delete", Method: "DELETE", Path: "/api/agent-network/providers/{providerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.ingress.ports.list":          {Name: "peers.ingress.ports.list", Method: "GET", Path: "/api/peers/{peerId}/ingress/ports", Mutation: ReadOnly},
	"peers.ingress.ports.get":           {Name: "peers.ingress.ports.get", Method: "GET", Path: "/api/peers/{peerId}/ingress/ports/{allocationId}", Mutation: ReadOnly},
	"peers.ingress.ports.create":        {Name: "peers.ingress.ports.create", Method: "POST", Path: "/api/peers/{peerId}/ingress/ports", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.ingress.ports.update":        {Name: "peers.ingress.ports.update", Method: "PUT", Path: "/api/peers/{peerId}/ingress/ports/{allocationId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.ingress.ports.delete":        {Name: "peers.ingress.ports.delete", Method: "DELETE", Path: "/api/peers/{peerId}/ingress/ports/{allocationId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.accessible":                  {Name: "peers.accessible", Method: "GET", Path: "/api/peers/{peerId}/accessible-peers", Mutation: ReadOnly},
	"networks.list":                     {Name: "networks.list", Method: "GET", Path: "/api/networks", Mutation: ReadOnly},
	"networks.get":                      {Name: "networks.get", Method: "GET", Path: "/api/networks/{networkId}", Mutation: ReadOnly},
	"networks.create":                   {Name: "networks.create", Method: "POST", Path: "/api/networks", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"routes.list":                       {Name: "routes.list", Method: "GET", Path: "/api/routes", Mutation: ReadOnly},
	"routes.get":                        {Name: "routes.get", Method: "GET", Path: "/api/routes/{routeId}", Mutation: ReadOnly},
	"users.current":                     {Name: "users.current", Method: "GET", Path: "/api/users/current", Mutation: ReadOnly},
	"users.list":                        {Name: "users.list", Method: "GET", Path: "/api/users", Mutation: ReadOnly},
	"users.invites":                     {Name: "users.invites", Method: "GET", Path: "/api/users/invites", Mutation: ReadOnly},
	"users.invites.create":              {Name: "users.invites.create", Method: "POST", Path: "/api/users/invites", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.invites.delete":              {Name: "users.invites.delete", Method: "DELETE", Path: "/api/users/invites/{inviteId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.invites.regenerate":          {Name: "users.invites.regenerate", Method: "POST", Path: "/api/users/invites/{inviteId}/regenerate", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"users.invites.accept":              {Name: "users.invites.accept", Method: "POST", Path: "/api/users/invites/{token}/accept", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: false, DispatcherAdmitted: true},
	"groups.list":                       {Name: "groups.list", Method: "GET", Path: "/api/groups", Mutation: ReadOnly},
	"groups.get":                        {Name: "groups.get", Method: "GET", Path: "/api/groups/{groupId}", Mutation: ReadOnly},
	"groups.create":                     {Name: "groups.create", Method: "POST", Path: "/api/groups", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"groups.update":                     {Name: "groups.update", Method: "PUT", Path: "/api/groups/{groupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"groups.delete":                     {Name: "groups.delete", Method: "DELETE", Path: "/api/groups/{groupId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.list":                        {Name: "peers.list", Method: "GET", Path: "/api/peers", Mutation: ReadOnly},
	"peers.get":                         {Name: "peers.get", Method: "GET", Path: "/api/peers/{peerId}", Mutation: ReadOnly},
	"policies.list":                     {Name: "policies.list", Method: "GET", Path: "/api/policies", Mutation: ReadOnly},
	"policies.get":                      {Name: "policies.get", Method: "GET", Path: "/api/policies/{policyId}", Mutation: ReadOnly},
	"policies.create":                   {Name: "policies.create", Method: "POST", Path: "/api/policies", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"policies.update":                   {Name: "policies.update", Method: "PUT", Path: "/api/policies/{policyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"policies.delete":                   {Name: "policies.delete", Method: "DELETE", Path: "/api/policies/{policyId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"routes.create":                     {Name: "routes.create", Method: "POST", Path: "/api/routes", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"routes.update":                     {Name: "routes.update", Method: "PUT", Path: "/api/routes/{routeId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"routes.delete":                     {Name: "routes.delete", Method: "DELETE", Path: "/api/routes/{routeId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.update":                      {Name: "peers.update", Method: "PUT", Path: "/api/peers/{peerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.delete":                      {Name: "peers.delete", Method: "DELETE", Path: "/api/peers/{peerId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.temporary_access.create":     {Name: "peers.temporary_access.create", Method: "POST", Path: "/api/peers/{peerId}/temporary-access", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: false, DispatcherAdmitted: true},
	"peers.edr.bypassed.list":           {Name: "peers.edr.bypassed.list", Method: "GET", Path: "/api/peers/edr/bypassed", Mutation: ReadOnly},
	"peers.edr.bypass.create":           {Name: "peers.edr.bypass.create", Method: "POST", Path: "/api/peers/{peer-id}/edr/bypass", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"peers.edr.bypass.delete":           {Name: "peers.edr.bypass.delete", Method: "DELETE", Path: "/api/peers/{peer-id}/edr/bypass", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.update":                   {Name: "networks.update", Method: "PUT", Path: "/api/networks/{networkId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
	"networks.delete":                   {Name: "networks.delete", Method: "DELETE", Path: "/api/networks/{networkId}", Mutation: Consequential, RequiresPreimage: true, RequiresReadBack: true, DispatcherAdmitted: true},
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
