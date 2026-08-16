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
nb accounts list --json
nb users list --json
nb users invites --json
nb dns nameservers list --json
nb dns settings --json
nb dns zones list --json
nb identity-providers list --json
nb posture-checks list --json
nb events --json
nb setup-keys list --json
nb locations countries --json
nb networks resources list <network-id> --json
nb networks routers list <network-id> --json
nb users tokens list <user-id> --json
```

Inventory commands are bounded reads. User and invite results intentionally omit
passwords and invite tokens, setup-key inventory omits the upstream setup-key
secret, and user-token inventory omits token values. These values are never part
of the stable `nb` output contract.

Request `--json` explicitly when consuming output as an agent. Consequential
NetBird changes belong under `nb stage`; do not invent a direct-write command.
Inspect the stage, acknowledge the exact findings, and apply an exact
`<stage-id>@<revision>` only after the caller's authority policy permits it.
