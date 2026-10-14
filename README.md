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

## License

Code in this repository is licensed under the [Apache License 2.0](LICENSE); documentation and specification files are licensed under [CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/).
