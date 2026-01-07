# Architecture Overview

`journal-cli` follows a **Pragmatic Clean Architecture**. The codebase is organized to separate domain logic from infrastructure concerns, though it favors simplicity over strict interface abstraction where appropriate for a CLI.

## Directory Structure

### `cmd/`
Entry points for the application.
-   `cmd/journal/main.go`: The main entry point. Reads flags and hands control to the `app` package.

### `internal/`
Private application code.

-   **`domain/`**: **(Inner Layer)**
    -   Contains pure data structures like `JournalEntry` and `Todo`.
    -   Has no dependencies on other internal packages.

-   **`app/`**: **(Application Layer)**
    -   Contains the "glue" code.
    -   Orchestrates the program flow: Load Config -> Load Data -> Run TUI -> Save Data.
    -   Currently implements the primary "Use Cases" implicitly.

-   **Infrastructure & Adapters**:
    -   `config/`: Handles loading and parsing `config.yaml`.
    -   `fs/`: File system helpers (wrappers for `os` and `io` with error handling).
    -   `markdown/`: Parsing and generation of Keep-a-Changelog style markdown journals.
    -   `tui/`: The Bubble Tea model and view logic. Handles the UI rendering and state machine.
    -   `updater/`: Handles the `self-update` mechanics (GitHub Releases API, binary replacement).
    -   `help/`: Structured help documentation logic.

## Data Flow

1.  **Startup**: `main.go` parses flags.
2.  **Orchestration**: `app.Run()` initializes config, template loader, and checks for existing files.
3.  **Interaction**: Control is handed to `tui.NewModel()`. The Bubble Tea runtime manages the event loop (keystrokes -> update model -> view).
4.  **Completion**: On exit, `app` retrieves the final state from the TUI model.
5.  **Persistence**: `app` calls `markdown.GenerateMarkdown()` and then `fs.WriteFile()` to save logic to disk.

## Design Decisions

See [ADR 001: Pragmatic Clean Architecture](adr/001-pragmatic-clean-architecture.md).
