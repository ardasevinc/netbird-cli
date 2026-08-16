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
```

Request `--json` explicitly when consuming output as an agent. Consequential
NetBird changes belong under `nb stage`; do not invent a direct-write command.
Inspect the stage, acknowledge the exact findings, and apply an exact
`<stage-id>@<revision>` only after the caller's authority policy permits it.
