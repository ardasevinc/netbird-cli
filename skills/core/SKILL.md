---
name: nb-core
description: Discover the installed nb command, schema, capability, and safety contracts.
---

Account updates are blocking and require exact preimage and read-back proof.
Account deletion is blocking and requires exact preimage and absence proof.
Posture-check creation is blocking and requires collection preimage and read-back proof.
Posture-check updates are blocking and require exact preimage and read-back proof.
Posture-check deletion is blocking and requires exact preimage and absence proof.
Ingress-peer creation is blocking and requires collection preimage and read-back proof.
Ingress-peer updates are blocking and require exact preimage and read-back proof.
Ingress-peer deletion is blocking and requires exact preimage and absence proof.
Peer ingress-port allocation creation, updates, and deletion are capability-gated, blocking, and require peer-scoped exact preimage plus read-back or absence proof because they can expose or remove external service mappings.
EDR compliance bypass and revoke are capability-gated, blocking security-control mutations; bypass grants immediate peer access outside normal compliance gating and requires exact collection proof.
Agent-network settings updates are capability-gated, blocking, and require exact preimage and read-back proof.
Agent-network settings creation is capability-gated, blocking, and requires exact preimage and read-back proof.
Agent-network settings deletion is capability-gated, blocking, and requires exact preimage and absence proof.
Agent-network budget-rule mutations are capability-gated, blocking, and require exact preimage plus read-back or absence proof.
Agent-network guardrail mutations are capability-gated, blocking, and require exact preimage plus read-back or absence proof.
Agent-network policy mutations are capability-gated, blocking, and require exact preimage plus read-back or absence proof.
Agent-network provider mutations are capability-gated, blocking, and require exact preimage plus read-back or absence proof. Provider API keys never persist in stages; use an external `api_key_ref` for one-time resolution at apply.
User create/update/delete and approval/rejection mutations are capability-gated, blocking, and require exact preimage plus read-back or absence proof.
Personal access token deletion is capability-gated, blocking, and requires exact token preimage plus absence proof; token values are never represented by the owned metadata model.
Personal access token creation is capability-gated and blocking; the token value is returned once in the successful apply result and is omitted from stages, receipts, logs, and errors.
Setup-key updates and deletion are capability-gated and blocking, with enrollment-impact evidence, exact preimage verification, and read-back or absence proof; setup-key creation remains one-time-secret gated.

Billing stages are Cloud-entitlement gated and blocking. AWS Marketplace actions
and subscription changes require an exact subscription preimage and read-back;
checkout creates an external payment session and returns its response proof only.

Temporary-access peer creation is capability-gated and blocking, with scoped-access impact evidence, exact target-peer preimage verification, and response-based proof because the ephemeral peer has no durable read-back surface.

Event-streaming integration create, update, and delete are capability-gated and blocking. Plans carry only `config_ref`; resolved receiver configuration exists in memory for dispatch, while masked server metadata is used for read-back.

Identity-provider create, update, and delete are capability-gated and blocking. Plans carry only `client_secret_ref`; the client secret is resolved in memory for dispatch and never persisted.

Reverse-proxy token creation and revocation are capability-gated and blocking. Creation returns the clear proxy token once, while revocation proves the token is absent from the server collection.
Authenticated invite create, delete, and regenerate mutations are capability-gated and blocking; invite tokens are returned once only and never persisted in stages, receipts, logs, or errors.
Password changes are capability-gated and blocking; staged requests carry only external `old_password_ref` and `new_password_ref` values, resolved in memory immediately before dispatch.
Resending a user invitation is capability-gated and blocking; the empty API response is proven by exact user preimage validation and unchanged user metadata read-back, with no token persisted or emitted.
Public invite acceptance is also staged and blocking; `invite_token_ref` and `password_ref` are resolved only for the unauthenticated acceptance request, with endpoint success as the effect proof because accepted invites no longer have pending metadata to read back.
Setup-key creation returns the clear key once in the successful apply result and omits it from stages, receipts, logs, and errors.
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
Peer deletion is blocking and requires the same absence proof.
Policy deletion is blocking and requires read-back confirmation of absence.
Policy creation is blocking and requires collection preimage and read-back proof.
Group deletion is blocking and requires the same absence proof.
Group creation is blocking and requires collection preimage and read-back proof.
Route deletion is blocking and requires the same absence proof.
Route creation is blocking and requires collection preimage and read-back proof.
Network deletion is blocking and requires the same absence proof.
Network creation is blocking and requires collection preimage and read-back proof.
DNS zone creation is blocking and requires collection preimage and read-back proof.
DNS zone deletion is blocking and requires the same absence proof.
DNS zone updates are blocking and require exact preimage and read-back proof.
DNS record creation is blocking and requires zone preimage and read-back proof.
DNS record updates are blocking and require exact zone/record preimage and read-back proof.
DNS record deletion is blocking and requires the same absence proof.
Nameserver-group creation is blocking and requires collection preimage and read-back proof.
Nameserver-group updates are blocking and require exact preimage and read-back proof.
Nameserver-group deletion is blocking and requires the same absence proof.
DNS-settings updates are blocking and require exact preimage and read-back proof.
Network resource deletion is blocking and requires the same absence proof.
Network resource address, enablement, or group changes are blocking.
Network resource creation is blocking and requires collection preimage and read-back proof.
Network router deletion is blocking and requires the same absence proof.
Network router updates are blocking.
Network router creation is blocking and requires collection preimage and read-back proof.
