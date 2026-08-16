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
go run ./cmd/nb version --jsonl
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
go run ./cmd/nb networks resources list <network-id> --json
go run ./cmd/nb networks routers list <network-id> --json
go run ./cmd/nb networks routers list-all --json
go run ./cmd/nb ingress list --json
go run ./cmd/nb peers ingress-ports list <peer-id> --json
go run ./cmd/nb dns nameservers list --json
go run ./cmd/nb dns settings --json
go run ./cmd/nb dns zones list --json
go run ./cmd/nb dns zones records list <zone-id> --json
go run ./cmd/nb identity-providers list --json
go run ./cmd/nb posture-checks list --json
go run ./cmd/nb events --json
go run ./cmd/nb events network-traffic --page 1 --page-size 100 --json
go run ./cmd/nb events proxy --page 1 --page-size 50 --json
go run ./cmd/nb setup-keys list --json
go run ./cmd/nb locations countries --json
go run ./cmd/nb users tokens list <user-id> --json
go run ./cmd/nb peers list --json
go run ./cmd/nb peers get <peer-id> --json
go run ./cmd/nb peers accessible <peer-id> --json
go run ./cmd/nb analyze reachability <peer-id> --json
go run ./cmd/nb policies list --json
go run ./cmd/nb policies get <policy-id> --json
go run ./cmd/nb stage create --from-json
go run ./cmd/nb apply <stage-id>@<revision>
```

Use `--json` explicitly for machine consumption. Core commands never silently
switch output modes based on whether stdout is a TTY.

Use `--jsonl` for a bounded stream: each line is an independent
`nb/v1/stream-event` record, and successful finite commands end with exactly
one `complete` line. `--json` and `--jsonl` are explicit output modes, never
TTY heuristics.

Reachability analysis reports server-reported accessible peers first, then
policy-group intersections as explanatory evidence. Any reachable peer without
an enabled `accept` rule match remains listed as unexplained; the analysis does
not infer data-plane certainty from policy shape alone.

The network-traffic and proxy event commands expose one server page at a time.
Their JSON responses retain the server pagination totals and mark the result
`partial` until the requested page reaches the final page.

Consequential changes are created as immutable local stages. `apply` accepts
only an exact stage ID and revision, rechecks the live preimage, journals the
dispatch intent, performs one remote mutation request, and verifies the result by
reading the resource back. Never point an unreviewed stage at a production
profile. User inventory never emits upstream password or invite-token fields;
those values, along with setup-key and token secrets, are treated as one-time
or secret material even when returned by the management API.

Group stages persist deterministic mutation-impact evidence. Metadata-only group
changes are marked reachability-neutral; unsupported state changes are marked
unknown and become blocking until the exact impact finding is acknowledged.

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
route/network, and DNS reads,
and proves a staged group update through remote read-back. It never targets a
production profile.

Release tags use the checked-in `.goreleaser.yaml` contract. It builds macOS and
Linux archives with injected version provenance, SHA-256 checksums, archive
SBOMs, keyless Sigstore bundles, and a Homebrew tap formula. The release
workflow is intentionally tag-triggered and drafts the GitHub release; it needs
the eventual repository and `HOMEBREW_TAP_TOKEN` secret before publication.

## Status

The v1 objective is to ship complete declared NetBird management API coverage,
truthful capability and verification reporting, safe staged mutation semantics,
runtime agent skills, and reproducible release artifacts. See the project
notes for the locked contracts and release claim.
