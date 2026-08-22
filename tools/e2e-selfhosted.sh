#!/usr/bin/env bash
set -Eeuo pipefail

on_error() {
	status=$?
	echo "self-hosted e2e failed at line ${BASH_LINENO[0]}" >&2
	return "$status"
}
trap on_error ERR

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
upstream_repo="${NB_E2E_NETBIRD_REPO:-/Users/arda/.agents/repo-ref/netbird}"
upstream_tag="v0.77.0"
image="${NB_E2E_NETBIRD_IMAGE:-nb-netbird-cli-e2e:v0.77.0}"
client_image="${NB_E2E_NETBIRD_CLIENT_IMAGE:-nb-netbird-cli-client-e2e:v0.77.0}"
container="nb-cli-e2e-${RANDOM}-${RANDOM}"
source_client="${container}-source"
destination_client="${container}-destination"
work_dir="$(mktemp -d "${TMPDIR:-/tmp}/nb-cli-e2e.XXXXXX")"
fixture_id=""
user_id=""
pat=""
base_url=""

cleanup() {
	set +e
	if [[ -n "$fixture_id" && -n "$base_url" && -n "$pat" ]]; then
		curl -fsS -X DELETE -H "Authorization: Bearer $pat" "$base_url/api/groups/$fixture_id" >/dev/null 2>&1
	fi
	if [[ -n "$user_id" && -n "$base_url" && -n "$pat" ]]; then
		curl -fsS -X DELETE -H "Authorization: Bearer $pat" "$base_url/api/users/$user_id" >/dev/null 2>&1
	fi
	docker rm -f "$source_client" "$destination_client" >/dev/null 2>&1
	docker rm -f "$container" >/dev/null 2>&1
	find "$work_dir" -type f -delete >/dev/null 2>&1
	find "$work_dir" -depth -type d -empty -delete >/dev/null 2>&1
}
trap cleanup EXIT

command -v docker >/dev/null || { echo "docker is required for self-hosted e2e" >&2; exit 2; }
command -v curl >/dev/null || { echo "curl is required for self-hosted e2e" >&2; exit 2; }
command -v jq >/dev/null || { echo "jq is required for self-hosted e2e" >&2; exit 2; }
test -d "$upstream_repo" || { echo "set NB_E2E_NETBIRD_REPO to a NetBird source checkout" >&2; exit 2; }
git -C "$upstream_repo" cat-file -e "$upstream_tag^{commit}" || { echo "missing pinned NetBird tag $upstream_tag" >&2; exit 2; }

expected_openapi_sha="3470cede60e3263e715b42c37f0d4534efa5ac54ab81b0e87081ade26e6dbe83"
if command -v sha256sum >/dev/null; then
	actual_openapi_sha="$(git -C "$upstream_repo" show "$upstream_tag:shared/management/http/api/openapi.yml" | sha256sum | awk '{print $1}')"
else
	actual_openapi_sha="$(git -C "$upstream_repo" show "$upstream_tag:shared/management/http/api/openapi.yml" | shasum -a 256 | awk '{print $1}')"
fi
test "$actual_openapi_sha" = "$expected_openapi_sha" || { echo "pinned OpenAPI hash mismatch" >&2; exit 2; }

if [[ -z "${NB_E2E_NETBIRD_IMAGE:-}" || -z "${NB_E2E_NETBIRD_CLIENT_IMAGE:-}" ]]; then
	upstream_context="$work_dir/netbird"
	mkdir -p "$upstream_context"
	git -C "$upstream_repo" archive "$upstream_tag" | tar -xf - -C "$upstream_context"
fi
if [[ -z "${NB_E2E_NETBIRD_IMAGE:-}" ]]; then
	docker build --quiet -f "$upstream_context/combined/Dockerfile.multistage" -t "$image" "$upstream_context" >/dev/null
fi
if [[ -z "${NB_E2E_NETBIRD_CLIENT_IMAGE:-}" ]]; then
	docker build --quiet -f "$upstream_context/e2e/harness/Dockerfile.client" -t "$client_image" "$upstream_context" >/dev/null
fi

mkdir -p "$work_dir/data"
sed "s|__EXPOSED_ADDRESS__|http://127.0.0.1:0|" "$repo_root/tests/e2e/netbird-config.yaml" > "$work_dir/config.yaml"
docker run -d --name "$container" -p 127.0.0.1::8080 -v "$work_dir:/nb" \
	-e NB_SETUP_PAT_ENABLED=true -e NB_DISABLE_GEOLOCATION=true \
	"$image" --config /nb/config.yaml >/dev/null

host_port=""
for _ in $(seq 1 60); do
	host_port="$(docker port "$container" 8080/tcp 2>/dev/null | sed 's/.*://')"
	if [[ -n "$host_port" ]] && curl -fsS "http://127.0.0.1:$host_port/api/instance" >/dev/null 2>&1; then
		break
	fi
	host_port=""
	sleep 2
done
test -n "$host_port" || { docker logs "$container" >&2; exit 1; }
base_url="http://127.0.0.1:$host_port"

bootstrap_password="$work_dir/bootstrap-password"
printf '%s' 'Netbird-e2e-Passw0rd!' > "$bootstrap_password"
chmod 600 "$bootstrap_password"
setup_json="$("$repo_root/bin/nb" --json setup bootstrap --url "$base_url" --email admin@netbird.test --name 'NB E2E Admin' --password-ref "file:$bootstrap_password" --create-pat --pat-expire-in 1)"
pat="$(jq -er '.data.personal_access_token // empty' <<<"$setup_json")"
account_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/accounts" | jq -er '.[0].id')"

profile="$work_dir/profile.toml"
state="$work_dir/ledger.db"
printf '%s\n' \
	"[profiles.default]" \
	"url = \"$base_url\"" \
	'credential_ref = "env:NB_E2E_TOKEN"' \
	"server_identity = \"$base_url\"" \
	"account_id = \"$account_id\"" > "$profile"

nb() {
	NB_E2E_TOKEN="$pat" "$repo_root/bin/nb" --config "$profile" --state "$state" "$@"
}

accounts_json="$(nb --json accounts list)"
test "$(jq -r '.operation' <<<"$accounts_json")" = "accounts.list"
test "$(jq -r '.data.accounts[0].id' <<<"$accounts_json")" = "$account_id"
users_json="$(nb --json users list)"
test "$(jq -r '.operation' <<<"$users_json")" = "users.list"

# Prove the user mutation preimage/read-back path against the live pinned
# server. NetBird exposes user reads through the collection endpoint, not
# GET /api/users/{userId}. Create a disposable service user, update its role
# through nb stage/apply, then remove it directly during cleanup.
user_id="$(curl -fsS -X POST "$base_url/api/users" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data '{"email":"nb-e2e-user@netbird.test","name":"NB E2E User","role":"user","auto_groups":[],"is_service_user":true}' | jq -er '.id')"
user_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/users" | jq -ce --arg id "$user_id" 'map(select(.id == $id)) | .[0]')"
user_after="$(jq '.role = "admin"' <<<"$user_before")"
user_request="$(jq -n --arg id "$user_id" --argjson after "$user_after" '{id:$id,role:$after.role,auto_groups:$after.auto_groups,is_blocked:$after.is_blocked}')"
user_plan="$(jq -n --arg id "$user_id" --argjson request "$user_request" --argjson before "$user_before" --argjson after "$user_after" '{operation:"users.update",request:$request,before:$before,intended_after:$after,findings:[]}')"
user_stage_json="$(printf '%s' "$user_plan" | nb --json stage create --from-json)"
user_stage_id="$(jq -er '.data.stage_id' <<<"$user_stage_json")"
user_revision="$(jq -er '.data.revision' <<<"$user_stage_json")"
user_apply_json="$(nb --json apply "$user_stage_id@$user_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$user_apply_json")" = "confirmed_success"
user_readback="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/users" | jq -ce --arg id "$user_id" 'map(select(.id == $id)) | .[0]')"
test "$(jq -r '.role' <<<"$user_readback")" = "admin"

# NetBird v0.77.0 publishes ingress routes in OpenAPI but the disposable OSS
# management server does not register them. Keep the absence explicit and
# prove nb preserves the safe HTTP context instead of calling it generic
# rejection.
ingress_status="$(curl -sS -o "$work_dir/ingress-body" -w '%{http_code}' -H "Authorization: Bearer $pat" "$base_url/api/ingress/peers")"
test "$ingress_status" = "404"
if nb --json ingress list >"$work_dir/ingress-stdout" 2>"$work_dir/ingress-stderr"; then
	echo "expected nb ingress list to report the unavailable pinned endpoint" >&2
	exit 1
fi
grep -Fq 'HTTP 404 GET /api/ingress/peers' "$work_dir/ingress-stdout"
grep -Fq 'HTTP 404 GET /api/ingress/peers' "$work_dir/ingress-stderr"

groups_json="$(nb --json groups list)"
test "$(jq -r '.data.groups | length' <<<"$groups_json")" -ge 1
routes_json="$(nb --json routes list)"
test "$(jq -r '.operation' <<<"$routes_json")" = "routes.list"
networks_json="$(nb --json networks list)"
test "$(jq -r '.operation' <<<"$networks_json")" = "networks.list"

# The built-in All group is intentionally immutable. Create a disposable
# fixture through the server API, then exercise the consequential update only
# through nb stage/apply and verify its read-back.
fixture_id="$(curl -fsS -X POST "$base_url/api/groups" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data '{"name":"nb-e2e-fixture"}' | jq -er '.id')"
before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/groups/$fixture_id")"
after="$(jq '.name = "nb-e2e-fixture-renamed"' <<<"$before")"
plan="$(jq -n --arg id "$fixture_id" --argjson before "$before" --argjson after "$after" '{operation:"groups.update",request:{id:$id,name:"nb-e2e-fixture-renamed"},before:$before,intended_after:$after,findings:[]}')"
stage_json="$(printf '%s' "$plan" | nb --json stage create --from-json)"
stage_id="$(jq -er '.data.stage_id' <<<"$stage_json")"
revision="$(jq -er '.data.revision' <<<"$stage_json")"
apply_json="$(nb --json apply "$stage_id@$revision")"
test "$(jq -r '.data.state' <<<"$apply_json")" = "confirmed_success"
readback_json="$(nb --json groups get "$fixture_id")"
test "$(jq -r '.data.group.name' <<<"$readback_json")" = "nb-e2e-fixture-renamed"

all_group_id="$(jq -er '.data.groups[] | select(.name == "All") | .id' <<<"$groups_json")"
default_policy_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/policies" | jq -er '.[] | select(.name == "Default") | .id')"
curl -fsS -X DELETE -H "Authorization: Bearer $pat" "$base_url/api/policies/$default_policy_id" >/dev/null

# The management endpoint named accessible-peers exposes symmetric network-map
# adjacency and discards the per-peer firewall direction. Enroll two disposable
# clients into distinct groups and create one non-bidirectional TCP flow. The
# new commands must retain both facts: adjacency in both maps, outbound-only at
# the source, and inbound-only at the destination. No packet probe is performed
# here, so observations must remain unknown and empty.
source_group_id="$(curl -fsS -X POST "$base_url/api/groups" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data '{"name":"nb-e2e-source"}' | jq -er '.id')"
destination_group_id="$(curl -fsS -X POST "$base_url/api/groups" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data '{"name":"nb-e2e-destination"}' | jq -er '.id')"
source_setup="$(jq -n --arg gid "$source_group_id" '{name:"nb-e2e-source",type:"one-off",expires_in:86400,auto_groups:[$gid],usage_limit:1}')"
destination_setup="$(jq -n --arg gid "$destination_group_id" '{name:"nb-e2e-destination",type:"one-off",expires_in:86400,auto_groups:[$gid],usage_limit:1}')"
source_key="$(curl -fsS -X POST "$base_url/api/setup-keys" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data "$source_setup" | jq -er '.key')"
destination_key="$(curl -fsS -X POST "$base_url/api/setup-keys" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data "$destination_setup" | jq -er '.key')"
management_url="http://host.docker.internal:$host_port"
docker run -d --name "$source_client" --cap-add NET_ADMIN --cap-add SYS_ADMIN --cap-add SYS_RESOURCE --device /dev/net/tun \
	-e NB_SETUP_KEY="$source_key" -e NB_MANAGEMENT_URL="$management_url" -e NB_HOSTNAME=nb-e2e-source \
	"$client_image" >/dev/null
docker run -d --name "$destination_client" --cap-add NET_ADMIN --cap-add SYS_ADMIN --cap-add SYS_RESOURCE --device /dev/net/tun \
	-e NB_SETUP_KEY="$destination_key" -e NB_MANAGEMENT_URL="$management_url" -e NB_HOSTNAME=nb-e2e-destination \
	"$client_image" >/dev/null

source_peer_id=""
destination_peer_id=""
for _ in $(seq 1 60); do
	peers_live="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/peers")"
	source_peer_id="$(jq -r '.[] | select(.name == "nb-e2e-source") | .id' <<<"$peers_live")"
	destination_peer_id="$(jq -r '.[] | select(.name == "nb-e2e-destination") | .id' <<<"$peers_live")"
	if [[ -n "$source_peer_id" && -n "$destination_peer_id" ]]; then
		break
	fi
	sleep 2
done
if [[ -z "$source_peer_id" || -z "$destination_peer_id" ]]; then
	docker logs "$source_client" >&2
	docker logs "$destination_client" >&2
	echo "disposable peers did not enroll" >&2
	exit 1
fi

one_way_policy="$(jq -n --arg source "$source_group_id" --arg destination "$destination_group_id" '{name:"nb-e2e-one-way",enabled:true,rules:[{name:"source-to-destination-https",description:"directionality fixture",enabled:true,action:"accept",protocol:"tcp",bidirectional:false,sources:[$source],destinations:[$destination],ports:["443"],port_ranges:[]}]}')"
curl -fsS -X POST "$base_url/api/policies" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data "$one_way_policy" >/dev/null

source_map=""
destination_map=""
for _ in $(seq 1 60); do
	source_map="$(nb --json peers network-map "$source_peer_id")"
	destination_map="$(nb --json peers network-map "$destination_peer_id")"
	if jq -e --arg id "$destination_peer_id" '.data.peers | any(.id == $id)' <<<"$source_map" >/dev/null && \
		jq -e --arg id "$source_peer_id" '.data.peers | any(.id == $id)' <<<"$destination_map" >/dev/null; then
		break
	fi
	sleep 2
done
jq -e --arg id "$destination_peer_id" '.data.peers | any(.id == $id)' <<<"$source_map" >/dev/null
jq -e --arg id "$source_peer_id" '.data.peers | any(.id == $id)' <<<"$destination_map" >/dev/null

source_access="$(nb --json analyze access "$source_peer_id")"
destination_access="$(nb --json analyze access "$destination_peer_id")"
if ! jq -e --arg id "$destination_peer_id" '.data.relations[] | select(.peer.id == $id) | .network_map_adjacent and .has_outbound_flows and (.has_inbound_flows | not) and (.outbound_flows | any(.protocol == "tcp" and (.ports == ["443"]) and (.bidirectional | not)))' <<<"$source_access" >/dev/null; then
	jq --arg id "$destination_peer_id" '.data.relations[] | select(.peer.id == $id)' <<<"$source_access" >&2
	echo "source access relation is not outbound-only" >&2
	exit 1
fi
if ! jq -e --arg id "$source_peer_id" '.data.relations[] | select(.peer.id == $id) | .network_map_adjacent and (.has_outbound_flows | not) and .has_inbound_flows and (.inbound_flows | any(.protocol == "tcp" and (.ports == ["443"]) and (.bidirectional | not)))' <<<"$destination_access" >/dev/null; then
	jq --arg id "$source_peer_id" '.data.relations[] | select(.peer.id == $id)' <<<"$destination_access" >&2
	echo "destination access relation is not inbound-only" >&2
	exit 1
fi
jq -e '.data.observations == [] and .data.completeness.observations.state == "unknown"' <<<"$source_access" >/dev/null
jq -e '.data.observations == [] and .data.completeness.observations.state == "unknown"' <<<"$destination_access" >/dev/null

# DNS zone and record lifecycles exercise nested targets and the target-field
# stripping used by both create and update dispatch.
zones_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones")"
zone_request="$(jq -n --arg gid "$all_group_id" '{name:"nb-e2e-zone",domain:"nb-e2e.internal",enabled:true,enable_search_domain:false,distribution_groups:[$gid]}')"
zone_plan="$(jq -n --argjson request "$zone_request" --argjson before "$zones_before" '{operation:"dns.zones.create",request:$request,before:$before,intended_after:{name:"nb-e2e-zone",domain:"nb-e2e.internal",enabled:true,enable_search_domain:false},findings:[]}')"
zone_stage="$(printf '%s' "$zone_plan" | nb --json stage create --from-json)"
zone_stage_id="$(jq -er '.data.stage_id' <<<"$zone_stage")"
zone_revision="$(jq -er '.data.revision' <<<"$zone_stage")"
zone_apply="$(nb --json apply "$zone_stage_id@$zone_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$zone_apply")" = "confirmed_success"
zone_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones" | jq -er '.[] | select(.name == "nb-e2e-zone") | .id')"
records_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones/$zone_id/records")"
record_request='{"name":"www.nb-e2e.internal","type":"A","content":"192.0.2.10","ttl":60}'
record_plan="$(jq -n --arg zone_id "$zone_id" --argjson request "$record_request" --argjson before "$records_before" '{operation:"dns.records.create",request:($request + {zone_id:$zone_id}),before:$before,intended_after:{name:"www.nb-e2e.internal",type:"A",content:"192.0.2.10",ttl:60},findings:[]}')"
record_stage="$(printf '%s' "$record_plan" | nb --json stage create --from-json)"
record_stage_id="$(jq -er '.data.stage_id' <<<"$record_stage")"
record_revision="$(jq -er '.data.revision' <<<"$record_stage")"
record_apply="$(nb --json apply "$record_stage_id@$record_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$record_apply")" = "confirmed_success"
record_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones/$zone_id/records" | jq -er '.[] | select(.name == "www.nb-e2e.internal") | .id')"
record_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones/$zone_id/records/$record_id")"
record_update_request="$(jq '{id:.id,zone_id:"'$zone_id'",name:.name,type:.type,content:"192.0.2.11",ttl:120}' <<<"$record_before")"
record_update_after="$(jq '{id:.id,name:.name,type:.type,content:"192.0.2.11",ttl:120}' <<<"$record_before")"
record_update_plan="$(jq -n --argjson request "$record_update_request" --argjson before "$record_before" --argjson after "$record_update_after" '{operation:"dns.records.update",request:$request,before:$before,intended_after:$after,findings:[]}')"
record_update_stage="$(printf '%s' "$record_update_plan" | nb --json stage create --from-json)"
record_update_stage_id="$(jq -er '.data.stage_id' <<<"$record_update_stage")"
record_update_revision="$(jq -er '.data.revision' <<<"$record_update_stage")"
record_update_apply="$(nb --json apply "$record_update_stage_id@$record_update_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$record_update_apply")" = "confirmed_success"
record_delete_plan="$(jq -n --arg zone_id "$zone_id" --arg record_id "$record_id" --argjson before "$record_update_after" '{operation:"dns.records.delete",request:{id:$record_id,zone_id:$zone_id},before:$before,intended_after:{},findings:[]}')"
record_delete_stage="$(printf '%s' "$record_delete_plan" | nb --json stage create --from-json)"
record_delete_stage_id="$(jq -er '.data.stage_id' <<<"$record_delete_stage")"
record_delete_revision="$(jq -er '.data.revision' <<<"$record_delete_stage")"
record_delete_apply="$(nb --json apply "$record_delete_stage_id@$record_delete_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$record_delete_apply")" = "confirmed_success"
zone_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/zones/$zone_id")"
zone_delete_plan="$(jq -n --arg id "$zone_id" --argjson before "$zone_before" '{operation:"dns.zones.delete",request:{id:$id},before:$before,intended_after:{},findings:[]}')"
zone_delete_stage="$(printf '%s' "$zone_delete_plan" | nb --json stage create --from-json)"
zone_delete_stage_id="$(jq -er '.data.stage_id' <<<"$zone_delete_stage")"
zone_delete_revision="$(jq -er '.data.revision' <<<"$zone_delete_stage")"
zone_delete_apply="$(nb --json apply "$zone_delete_stage_id@$zone_delete_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$zone_delete_apply")" = "confirmed_success"

settings_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/settings")"
settings_on="$(jq -n --arg gid "$all_group_id" '{disabled_management_groups:[$gid]}')"
settings_plan="$(jq -n --argjson request "$settings_on" --argjson before "$settings_before" '{operation:"dns.settings.update",request:$request,before:$before,intended_after:$request,findings:[]}')"
settings_stage="$(printf '%s' "$settings_plan" | nb --json stage create --from-json)"
settings_stage_id="$(jq -er '.data.stage_id' <<<"$settings_stage")"
settings_revision="$(jq -er '.data.revision' <<<"$settings_stage")"
settings_apply="$(nb --json apply "$settings_stage_id@$settings_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$settings_apply")" = "confirmed_success"
settings_before_restore="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/dns/settings")"
settings_off='{"disabled_management_groups":[]}'
settings_restore_plan="$(jq -n --argjson request "$settings_off" --argjson before "$settings_before_restore" '{operation:"dns.settings.update",request:$request,before:$before,intended_after:$request,findings:[]}')"
settings_restore_stage="$(printf '%s' "$settings_restore_plan" | nb --json stage create --from-json)"
settings_restore_stage_id="$(jq -er '.data.stage_id' <<<"$settings_restore_stage")"
settings_restore_revision="$(jq -er '.data.revision' <<<"$settings_restore_stage")"
settings_restore_apply="$(nb --json apply "$settings_restore_stage_id@$settings_restore_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$settings_restore_apply")" = "confirmed_success"

# Policy and posture-check payloads use different shapes for writes and reads:
# policy writes use group IDs, while reads expand groups into objects. Keep the
# staged intent to fields that are stable across both representations.
policies_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/policies")"
policy_request="$(jq -n --arg gid "$all_group_id" '{name:"nb-e2e-policy",enabled:true,rules:[{name:"allow-all",description:"nb e2e",enabled:true,action:"accept",protocol:"all",bidirectional:true,sources:[$gid],destinations:[$gid],ports:[],port_ranges:[]}] }')"
policy_plan="$(jq -n --argjson request "$policy_request" --argjson before "$policies_before" '{operation:"policies.create",request:$request,before:$before,intended_after:{name:"nb-e2e-policy",enabled:true},findings:[]}')"
policy_stage="$(printf '%s' "$policy_plan" | nb --json stage create --from-json)"
policy_stage_id="$(jq -er '.data.stage_id' <<<"$policy_stage")"
policy_revision="$(jq -er '.data.revision' <<<"$policy_stage")"
policy_apply="$(nb --json apply "$policy_stage_id@$policy_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$policy_apply")" = "confirmed_success"
policy_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/policies" | jq -er '.[] | select(.name == "nb-e2e-policy") | .id')"
policy_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/policies/$policy_id")"
policy_delete_plan="$(jq -n --arg id "$policy_id" --argjson before "$policy_before" '{operation:"policies.delete",request:{id:$id},before:$before,intended_after:{},findings:[]}')"
policy_delete_stage="$(printf '%s' "$policy_delete_plan" | nb --json stage create --from-json)"
policy_delete_id="$(jq -er '.data.stage_id' <<<"$policy_delete_stage")"
policy_delete_revision="$(jq -er '.data.revision' <<<"$policy_delete_stage")"
policy_delete_apply="$(nb --json apply "$policy_delete_id@$policy_delete_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$policy_delete_apply")" = "confirmed_success"
test -z "$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/policies" | jq -r '.[] | select(.name == "nb-e2e-policy") | .id')"

posture_before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/posture-checks")"
posture_request='{"name":"nb-e2e-posture","checks":{"process_check":{"processes":[{"linux_path":"/usr/bin/netbird"}]}}}'
posture_plan="$(jq -n --argjson request "$posture_request" --argjson before "$posture_before" '{operation:"posture_checks.create",request:$request,before:$before,intended_after:{name:"nb-e2e-posture"},findings:[]}')"
posture_stage="$(printf '%s' "$posture_plan" | nb --json stage create --from-json)"
posture_stage_id="$(jq -er '.data.stage_id' <<<"$posture_stage")"
posture_revision="$(jq -er '.data.revision' <<<"$posture_stage")"
posture_apply="$(nb --json apply "$posture_stage_id@$posture_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$posture_apply")" = "confirmed_success"
posture_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/posture-checks" | jq -er '.[] | select(.name == "nb-e2e-posture") | .id')"
posture_before_one="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/posture-checks/$posture_id")"
posture_delete_plan="$(jq -n --arg id "$posture_id" --argjson before "$posture_before_one" '{operation:"posture_checks.delete",request:{id:$id},before:$before,intended_after:{},findings:[]}')"
posture_delete_stage="$(printf '%s' "$posture_delete_plan" | nb --json stage create --from-json)"
posture_delete_id="$(jq -er '.data.stage_id' <<<"$posture_delete_stage")"
posture_delete_revision="$(jq -er '.data.revision' <<<"$posture_delete_stage")"
posture_delete_apply="$(nb --json apply "$posture_delete_id@$posture_delete_revision" --ack-all-blocking)"
test "$(jq -r '.data.state' <<<"$posture_delete_apply")" = "confirmed_success"
test -z "$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/posture-checks" | jq -r '.[] | select(.name == "nb-e2e-posture") | .id')"

printf 'self-hosted e2e passed: pinned %s, symmetric map adjacency, asymmetric configured flows, staged mutations, read-back confirmed\n' "$upstream_tag"
