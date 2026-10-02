---
format: https://specscore.md/idea-specification
status: Draft
---

# Idea: DataTug desktop app (macOS first)

**Status:** Draft
**Date:** 2026-10-02
**Owner:** alex
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** —

## Problem Statement

How might we give people who already rely on DataTug a native desktop app that reaches local databases without the browser-to-localhost bridge, without adding an install step in front of first value?

## Context

Raised by the founder on 2026-10-02 while planning DataTug.ai ("Should we also create a mac hybrid app? Using capacitor?"). Founder ruling the same day: *"I think CLI + web app is good enough for now and maybe even better. But record a desktop app idea and create a GitHub issue."* So this Idea is parked, not scheduled.

Today DataTug is the CLI/TUI (`datatug-cli`) plus the hosted web app at datatug.app (`datatug-apps`, Angular + Ionic). The web app reaches local data through a loopback bridge: `datatug chat` / `datatug serve` open `https://datatug.app<path>#h=127.0.0.1:<port>&t=<token>` and the page connects back over a WebSocket. That bridge depends on browsers continuing to let a public HTTPS page call a local address, which browsers have been restricting. The web app also keeps bring-your-own AI keys in browser `localStorage`.

The first-touch advantage of DataTug is value before install: the in-browser demo. A desktop app must not be put in front of that.

## Recommended Direction

Keep CLI + web as the product for now. When a desktop app is picked up, wrap the existing `datatug-apps` web build in a desktop shell that ships the Go `datatug` binary as a supervised sidecar, so the app reaches local databases with no separate `datatug serve` and no browser-to-localhost bridge, and stores credentials in the OS keychain.

Agent recommendation (not founder-ruled): use Tauri for the shell rather than Capacitor. Capacitor's supported platforms are iOS, Android and web; reaching macOS means Mac Catalyst (sandboxed, awkward for a sidecar binary) or the community Electron platform. Tauri treats sidecar binaries, keychain access, signing and auto-update as first-class, produces a small download, and the same project later builds for Windows and Linux. Capacitor remains the natural choice if iOS/Android apps are wanted, since the Ionic app is already the right base.

Position it as the upgrade for people who already use DataTug, never as the entry point.

## Alternatives Considered

- **Capacitor via Mac Catalyst.** Reuses a future iOS build, but the sandbox makes bundling and launching the Go binary awkward, and the binary is the point of a desktop build.
- **Capacitor community Electron platform.** Works, but community-maintained, a large download, and no advantage over Tauri for a sidecar.
- **Do nothing (CLI + web only).** The current ruling. Loses only if browser restrictions on public-page-to-localhost calls break the bridge, or if users ask for keychain-stored credentials.

## MVP Scope

A timeboxed spike (about two days): wrap the current `datatug-apps` build in a Tauri shell on macOS with the `datatug` binary as a sidecar, open the Chinook demo project, and run one chat query against a local database with no `datatug serve` started by hand and no `#h=` bridge. Signed and notarised build not required for the spike.

## Not Doing (and Why)

- A desktop-only feature set. One web codebase; the shell adds local access and keychain storage only.
- Windows and Linux builds in the MVP. macOS first; the others follow from the same project.
- Mobile apps. Separate question, Capacitor-based if pursued.
- Requiring install before the demo. The browser demo stays the first touch.

## Key Assumptions to Validate

| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | The existing web build runs unchanged inside a desktop webview with the CLI as a sidecar | The MVP spike |
| Must-be-true | Users want local access without a terminal step enough to install an app | Ask active CLI + web users; count bridge-related support issues |
| Should-be-true | Browser restrictions on public-page-to-localhost calls will make the bridge unreliable | Test the bridge in current Chrome, Safari and Firefox; track failures |
| Might-be-true | A desktop listing (Homebrew cask, direct download) is a meaningful acquisition channel | Compare install sources after a first release |


## SpecScore Integration

- **New Features this would create:** TBD at design time
- **Existing Features affected:** surfaces (a new surface), the CLI/web chat bridge
- **Dependencies:** `datatug-apps` web build, `datatug-cli` binary, macOS signing and notarisation

## Open Questions

- Tauri or another shell: confirm in the spike (agent recommendation is Tauri).
- Which trigger un-parks this: bridge breakage, user demand, or a distribution goal?
