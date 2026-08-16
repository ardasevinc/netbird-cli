# nb

`nb` is an unofficial, agent-first management CLI for NetBird. It complements
NetBird's peer-oriented official CLI with bounded inspection, explainable
analysis, staged consequential changes, and durable evidence.

The project is in active v1 implementation. The intended public contract is
documented in the companion project handoff and will become executable as the
vertical slices land.

## Current command spine

```sh
go run ./cmd/nb version
go run ./cmd/nb schema list
go run ./cmd/nb skills list
go run ./cmd/nb coverage
go run ./cmd/nb groups list --json
go run ./cmd/nb groups get <group-id> --json
go run ./cmd/nb accounts list --json
go run ./cmd/nb users list --json
go run ./cmd/nb users invites --json
go run ./cmd/nb routes list --json
go run ./cmd/nb networks list --json
go run ./cmd/nb peers list --json
go run ./cmd/nb peers get <peer-id> --json
go run ./cmd/nb policies list --json
go run ./cmd/nb policies get <policy-id> --json
go run ./cmd/nb stage create --from-json
go run ./cmd/nb apply <stage-id>@<revision>
```

Use `--json` explicitly for machine consumption. Core commands never silently
switch output modes based on whether stdout is a TTY.

Consequential changes are created as immutable local stages. `apply` accepts
only an exact stage ID and revision, rechecks the live preimage, journals the
dispatch intent, performs one remote mutation request, and verifies the result by
reading the resource back. Never point an unreviewed stage at a production
profile. User inventory never emits upstream password or invite-token fields;
those values are treated as one-time or secret material even when returned by
the management API.

## Development

```sh
just gate
# requires Docker, jq, and NB_E2E_NETBIRD_REPO (or the local reference checkout)
NB_E2E_NETBIRD_IMAGE=nb-netbird-cli-e2e:v0.77.0 just e2e-selfhosted
```

The local gate covers formatting, tests, race detection, vet, static analysis,
security checks, module verification, license boundaries, generated schemas,
cross-builds, and diff cleanliness.

The opt-in self-hosted lane pins NetBird `v0.77.0`, boots a disposable combined
server, bootstraps a throwaway admin PAT, exercises account/user/group,
route/network reads,
and proves a staged group update through remote read-back. It never targets a
production profile.

## Status

The v1 objective is to ship complete declared NetBird management API coverage,
truthful capability and verification reporting, safe staged mutation semantics,
runtime agent skills, and reproducible release artifacts. See the project
notes for the locked contracts and release claim.
