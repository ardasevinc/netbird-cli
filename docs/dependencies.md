# Dependency decisions

Dependencies are admitted narrowly and reviewed against the repository's
current Go CLI patterns. Versions and evidence are refreshed when a dependency
changes.

| Module | Version | Use | Evidence checked |
| --- | --- | --- | --- |
| `github.com/spf13/cobra` | `v1.10.2` | command parsing | already used by `tele` and `mattermost-cli`; active upstream release line |
| `github.com/pelletier/go-toml/v2` | `v2.4.3` | named profile decoding | already used by `tele` and `mattermost-cli`; focused TOML implementation |
| `modernc.org/sqlite` | `v1.53.0` | pure-Go local mutation ledger | already used by `mattermost-cli`; avoids CGO in the supported CLI builds |

The standard library remains preferred for HTTP, TLS, JSON, hashing, process
boundaries, and test servers. NetBird source dependencies are admitted only
after the coverage and mixed-license checks classify their package families.

## Changelog

- 2026-08-16 — codex — recorded the initial dependency set and local evidence
