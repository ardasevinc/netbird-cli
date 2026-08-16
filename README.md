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
```

Use `--json` explicitly for machine consumption. Core commands never silently
switch output modes based on whether stdout is a TTY.

## Development

```sh
just gate
```

The local gate covers formatting, tests, race detection, vet, static analysis,
security checks, module verification, license boundaries, generated schemas,
cross-builds, and diff cleanliness.

## Status

The v1 objective is to ship complete declared NetBird management API coverage,
truthful capability and verification reporting, safe staged mutation semantics,
runtime agent skills, and reproducible release artifacts. See the project
notes for the locked contracts and release claim.
