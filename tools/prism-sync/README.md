# prism-sync

Copies canonical templates and rubrics from
[grokify/prism-roadmap](https://github.com/grokify/prism-roadmap) into
this repo's embedded workflow catalog (`pkg/workflows/default/`). This is a
physical copy, not a Go import: the catalog is `//go:embed`-based
(filesystem-free loading at runtime), so the files have to actually live in
this repo's tree — a versioned dependency alone can't put markdown/YAML
where `go:embed` needs it.

It's a separate nested module (its own `go.mod`) specifically so
`prism-roadmap`'s dependency tree never enters the main module graph — this
is a one-time-per-run sync tool, not a runtime dependency of anything else
in this repo.

See `manifest.yaml` for exactly which templates/rubrics are synced from
prism-roadmap and which workflow families they land in — currently the
`enterprise` (bmc, opportunity-spec), `continuous-discovery` (ost,
assumption-map, discovery-snapshot), and the entire `v2mom` family (7
files).

## Current status: manual only, not wired into CI

There's an active local `replace` directive in this module's `go.mod`:

```
replace github.com/grokify/prism-roadmap => ../../../../grokify/prism-roadmap
```

This exists because the v2mom/ost templates and rubrics aren't in a tagged
prism-roadmap release yet (confirmed: `v0.16.1`, the latest tag as of
2026-08-15, doesn't have `templates/v2mom-summary.md` — it only exists on
an untagged commit past that tag). Until prism-roadmap ships a release
containing this content, the tool **cannot run from a fresh clone or in
CI** — only a developer with a local `prism-roadmap` checkout at that exact
relative path can run it.

**Deliberately not wired into a scheduled CI job right now.** A cron-triggered
`-check` run would fail on every single invocation until the replace is
resolved, training people to ignore a permanently-red check — worse than no
automation. Automate this once it can actually pass, not before.

### Running it manually (until then)

Requires a local checkout of `grokify/prism-roadmap` at
`../../../../grokify/prism-roadmap` relative to this directory (i.e., a
sibling of this org's checkouts under the same `github.com/` root).

```bash
# Check for drift without writing anything
go run -C tools/prism-sync . -check

# Actually sync (writes updated copies + provenance headers)
go run -C tools/prism-sync .
```

Run this periodically (e.g. whenever someone's doing work that touches the
`enterprise`, `continuous-discovery`, or `v2mom` workflow families) to catch
upstream drift in the meantime.

### Graduating to automated CI

Once `grokify/prism-roadmap` tags a release containing the v2mom/ost
content:

1. Drop the `replace` directive in `tools/prism-sync/go.mod`; require the
   tagged version instead.
2. Confirm `go run -C tools/prism-sync . -check` passes from a fresh clone.
3. Add a scheduled GitHub Actions workflow (weekly cron) running the same
   `-check` command, failing the job on drift.

At that point this becomes real automated drift detection instead of a
manual reminder.
