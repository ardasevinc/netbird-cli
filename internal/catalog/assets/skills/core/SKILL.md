---
name: nb-core
description: Discover the installed nb command, schema, capability, and safety contracts.
---

# nb core

Use the installed binary as the source of truth for its current machine
contract. Start with:

```sh
nb skills list
nb skills get core
nb schema list
nb capabilities --json
nb version --jsonl
nb api get get.api.reverse_proxies.services --json
nb accounts list --json
nb users list --json
nb users invites --json
nb dns nameservers list --json
nb dns settings --json
nb dns zones list --json
nb dns zones records list <zone-id> --json
nb identity-providers list --json
nb posture-checks list --json
nb events --json
nb events network-traffic --page 1 --page-size 100 --json
nb events proxy --page 1 --page-size 50 --json
nb setup-keys list --json
nb locations countries --json
nb networks resources list <network-id> --json
nb networks routers list <network-id> --json
nb ingress list --json
nb peers ingress-ports list <peer-id> --json
nb peers accessible <peer-id> --json
nb analyze reachability <peer-id> --json
nb users tokens list <user-id> --json
nb stage create --from-json
```

Inventory commands are bounded reads. User and invite results intentionally omit
passwords and invite tokens, setup-key inventory omits the upstream setup-key
secret, and user-token inventory omits token values. These values are never part
of the stable `nb` output contract.

`nb analyze reachability` uses the server-reported accessible-peer inventory as
the authoritative reachability result and attaches policy-group intersections
as explanatory evidence. Reachable peers without an enabled accept-rule match
are surfaced as unexplained rather than silently classified.

Request `--json` explicitly when consuming output as an agent. Consequential
NetBird changes belong under `nb stage`; do not invent a direct-write command.
Inspect the stage, acknowledge the exact findings, and apply an exact
`<stage-id>@<revision>` only after the caller's authority policy permits it.
Inspect the stage impact evidence as part of that review. Metadata-only group
changes and policy metadata changes are reachability-neutral; unknown group
impact, policy rule changes, and route behavior changes are blocking findings.
Peer access or connectivity changes are also blocking findings.
Network topology changes are also blocking findings.
Policy deletion is blocking and requires read-back confirmation of absence.
Group deletion is blocking and requires the same absence proof.
Route deletion is blocking and requires the same absence proof.
Network deletion is blocking and requires the same absence proof.
