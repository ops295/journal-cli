# 1. Adoption of Pragmatic Clean Architecture

Date: 2026-01-06

## Status

Accepted

## Context

The `journal-cli` project aims to be a robust, maintainable CLI application. We want to structure the code in a way that separates concerns, making it testable and easy to extend. However, strict adherence to Clean Architecture (with full interface abstraction for every layer) can introduce unnecessary boilerplate for a CLI tool of this size.

## Decision

We will adopt a **Pragmatic Clean Architecture** style.

### Layers

1.  **Domain (`internal/domain`)**:
    -   Contains pure business entities (e.g., `JournalEntry`, `Todo`).
    -   No external dependencies.

2.  **Application (`internal/app`)**:
    -   Orchestrates the application flow (e.g., loading config, checking files, running the loop).
    -   Currently acts as both "Use Case" and "Controller".
    -   *Constraint*: Should prioritize business logic over UI specifics where possible.

3.  **Adapters/Infrastructure (`internal/fs`, `internal/markdown`, `internal/tui`, `internal/updater`)**:
    -   Implement specific details (FileSystem, Markdown parsing, Terminal UI, GitHub Updates).
    -   The Application layer calls these helpers directly.

## Consequences

### Positive
-   **Simplicity**: Less boilerplate than defining interfaces for every file operation or markdown parser.
-   **Speed**: Faster development iteration.
-   **Clarity**: Directory structure clearly indicates what each package does.

### Negative
-   **Coupling**: `internal/app` is coupled to concrete implementations of `fs` and `tui`.
-   **Testing**: We cannot easily mock the FileSystem or TUI in integration tests without refactoring `app` to use interfaces.

### Mitigation
-   We accept this coupling for now.
-   If a component (like `fs`) becomes complex or needs interchangeable backends (e.g., S3 support), we will refactor it behind an interface at that time.
