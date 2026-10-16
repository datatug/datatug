# DataTug

**DataTug** is a context-aware, semantic data exploration and troubleshooting environment across operational data sources. It turns distributed operational data into a navigable graph of meaning, relationships, and reusable queries — stored as Git-native project files and usable across the terminal, browser, and AI agents.

## Surfaces & Ecosystem

- **DataTug CLI** — Command-line interface, terminal UI (TUI), and local execution server (`datatug serve`, `datatug demo`).
- **[Web App (`datatug.app`)](https://datatug.app)** — Browser workspace for interactive exploration, portable queries ([TugQL](https://tugql.org)), and comparisons.
- **[`datatug-core`](https://github.com/datatug/datatug-core)** — Core Go engine, project model, and storage abstractions.
- **[`datatug-demo-projects`](https://github.com/datatug/datatug-demo-projects)** — Sample Git-native DataTug projects and datasets.

## Quick Start

Install the CLI via Homebrew (macOS / Linux) and launch the interactive demo:

```sh
brew install datatug/tap/datatug
datatug demo
```

## Issues & Feature Requests

This repository is the main community hub for **GitHub Issues**, bug reports, and feature discussions across DataTug:

- [Open an Issue or Feature Request](https://github.com/datatug/datatug/issues)
- [Public Specifications & Ideas](spec/README.md)

## Under the hood

DataTug is built through serious engineering: we build on open source and contribute tools and standards back. These articles share some of the work behind the product. The OVDB session storage design is still in progress.

<a href="https://datatug.app/en/under-the-hood/tools-we-built/"><img src="https://datatug.app/static/illustrations/dev-shipyard-768.webp" width="160" alt="Illustration of tools being built in a software workshop" /></a>

**[The tools we built to build DataTug](https://datatug.app/en/under-the-hood/tools-we-built/)** — A look at the developer tools we created to build, ship, and maintain DataTug.

<a href="https://datatug.app/en/under-the-hood/ovdb-session-storage/"><img src="https://datatug.app/static/illustrations/ovdb-context-desk-768.webp" width="160" alt="Illustration of a workspace for organizing session context" /></a>

**[Why DataTug uses OVDB for working context](https://datatug.app/en/under-the-hood/ovdb-session-storage/)** — The architecture in progress for using OVDB to keep working context available across sessions.

[Browse all Under the hood articles](https://datatug.app/en/under-the-hood/).

## License

Code in this repository is licensed under the [Apache License 2.0](LICENSE); documentation and specification files are licensed under [CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/).
