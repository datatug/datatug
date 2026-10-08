# OVDB Directory PR35 post-merge proof

Observed 2026-10-08 around 22:10–22:15 UTC. Runtime: Codex Sol on local host. Read-only audit; report file is the only written artifact. No ECB feed fetched, provider query initiated, UI control submitted, or repository edited.

## Source landing and cleanup

- GitHub PR https://github.com/openvaultdb/ovdb-directory/pull/35 has state MERGED, mergedAt `2026-10-08T22:10:19Z`, mergeCommit `c37ac63806cda9021c56daa86255dce1a4b18b9a`, headRefName `ecb-directory-direct-browser`.
- Authoritative GitHub `branches/main` API reports exactly `c37ac63806cda9021c56daa86255dce1a4b18b9a`, rechecked after deployment. Local `origin/main` and canonical main HEAD report the same SHA.
- Canonical worktree `openvaultdb/ovdb-directory canonical checkout` is the sole local worktree. `git worktree list --porcelain` reports main at the merge SHA; `git branch` lists only main; `git status --porcelain` emits nothing.
- GitHub `git/matching-refs/heads/ecb-directory-direct-browser` returns `[]`, proving the remote task branch is gone. A stale local remote-tracking `origin/ecb-directory-direct-browser` remains in `git branch -a`; that is not a live remote branch. It was left untouched under read-only scope. Other remote-tracking refs were not claimed as task state.

## CI and deployment

- PR35 test check passed: https://github.com/openvaultdb/ovdb-directory/actions/runs/37851070820/job/113563976732 (3m22s from `gh pr checks 35`).
- Canonical `origin/main:.github/workflows/ci.yml` runs on pull requests only. Main push validation runs in Deploy, preventing duplicate main check billing.
- Exact merge-SHA Deploy run: https://github.com/openvaultdb/ovdb-directory/actions/runs/37851758615. Final API status `completed`, conclusion `success`; job conclusion `success`.
- Workflow evidence: unit tests, offline source artifact acceptance, syntax/fixture/wrangler dry-run check, production build/guard, fixture build, Playwright installation, desktop/375px e2e, Deploy, Smoke check and post-job cleanup all succeeded. The no-credentials skip notice was skipped, so it was an actual deploy.
- Deployment log: `2026-10-08T22:13:29.0506243Z Uploaded ovdb-directory (3.74 sec)`; `2026-10-08T22:13:30.5650618Z Current Version ID: e960ca6e-dbec-4928-90df-aee26548bd19`.
- Smoke log: `2026-10-08T22:13:31.0482702Z Deployed and confirmed: https://directory.openvaultdb.com serves commit c37ac63806cd.`

## Deployed artifact and dormant ECB proof

A fresh HTTP GET of https://directory.openvaultdb.com/build-info.json returned:

```json
{
  "format": "ovdb-directory-build/1",
  "source": "url",
  "commit": "c37ac63806cda9021c56daa86255dce1a4b18b9a",
  "indexCommits": {
    "directory": "480f369bedf4c3c8c6d1d67c32a083c6c820de54",
    "modelspecRegistry": "4002c097e73aba0b41c247f2a50776a228ad662f"
  },
  "indexChecksums": {
    "directory": "sha256:092035ad05ed4e278c458783817af3dafbfec28bd8ca398018b3a4a1a58c2fd8",
    "modelspecRegistry": "sha256:c191e57a271f4ae75a60dd0f0380e4ca11295eeebb9f32f09582f700622d4377",
    "directorySources": "sha256:363f6825cebc9ef119523d59ed346798c6fddba466927cda6317c9697d9042dd"
  }
}
```

Canonical source read from `git show origin/main`, at the verified merge SHA:

- `src/render.mjs:586` renders `<button type="submit" data-ecb-execute disabled>Execute</button>`.
- `src/browser/ecb-entry.mjs` mounts the reader without any injected admission/configuration.
- `src/browser/ecb-reader.mjs` defaults to `createEcbSession()`, sets Execute disabled whenever `!session.configured()`, and refuses submit in that state.
- `src/browser/ecb-session.mjs` considers a session configured only when an `admitExecutor` function is supplied; the default supplies none.

Fresh public browser load and post-deployment reload of https://directory.openvaultdb.com/sources/ecb-daily/ visibly report:

- Inactive; Query activation blocked.
- Heading: `Controlled transient reader — unconfigured`.
- Updated PR35 copy: direct browser HTTP is the intended default; an OVDB proxy requires separate admission when the direct protocol is unavailable.
- Execute disabled; Cancel read disabled.
- Status: `Reader unconfigured. Rights admission and browser access remain pending. No provider request is made.`

No provider links or Execute were clicked. This proves the page's observed dormant state and source wiring; it is not a packet capture or a live ECB transport acceptance test.

## Verdict

Source landing, task branch/worktree cleanup, final main validation/deployment, and exact deployed artifact are independently verified. ECB execution remains dormant in canonical main and the observed public site. The live artifact initially showed previous SHA `3a789ea00d65e329c0ee64f51474529ac6c6e185` while deployment ran; the final checks above confirm the new SHA after Deploy and Smoke succeeded.

Falsifiers: a later remote main SHA, live build marker mismatch, changed admission wiring, or enabled Execute on a fresh public load would invalidate the corresponding claim. Receipts describe the observed point in time.
