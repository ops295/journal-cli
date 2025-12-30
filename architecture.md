# Clean Architecture in Journal CLI

This project follows the principles of **Clean Architecture** to ensure separation of concerns, testability, and maintainability. The code is organized into concentric layers, with dependencies pointing inwards.

## Architectural Layers

### 1. Entities (Domain Layer)
**Location**: `internal/domain`

This is the innermost layer. It contains the core business objects of the application. These entities are plain Go structs and have **no dependencies** on outer layers (no filesystem, no TUI, no config).

- **`JournalEntry`**: Represents a daily journal entry with mood, energy, todos, and answers.
- **`Todo`**: Represents a single task.

### 2. Use Cases (Application Layer)
**Location**: `internal/todo`, `internal/markdown`, `internal/app`

This layer contains application-specific business rules. It orchestrates the flow of data to and from the entities.

- **`internal/todo`**: Contains the logic for calculating the backlog (business rule: "unchecked todos from yesterday become today's backlog").
- **`internal/markdown`**: Handles the conversion between Domain Entities and the persistence format (Markdown). This acts as a data mapper.
- **`internal/app`**: The **Application Orchestrator**. It acts as the main entry point for the business logic, wiring together the config, templates, and TUI.

### 3. Interface Adapters (Presentation Layer)
**Location**: `internal/tui`, `cmd`

This layer converts data from the format most convenient for the use cases and entities, to the format most convenient for some external agency (in this case, the Terminal).

- **`internal/tui`**: The **Bubble Tea** model. It handles user input (keyboard events) and renders the UI. It interacts with the `JournalEntry` domain object but delegates persistence to the outer layers (via the App Orchestrator).
- **`cmd/journal`**: The entry point. It initializes the application and dependencies.

### 4. Frameworks & Drivers (Infrastructure Layer)
**Location**: `internal/fs`, `internal/config`, `internal/template`

This is the outermost layer. It contains details such as the filesystem, configuration files, and external tools.

- **`internal/fs`**: Low-level wrappers around `os` and `filepath` to handle file I/O.
- **`internal/config`**: Knows how to find and parse the `config.yaml` file from the OS-specific user config directory.
- **`internal/template`**: Knows how to scan the `templates` directory and parse YAML files into template structures.

## Dependency Rule
The source code dependencies only point inwards.
- `internal/domain` knows nothing about `internal/tui` or `internal/fs`.
- `internal/tui` imports `internal/domain` but doesn't know about the specific filesystem implementation details (ideally, though in this simple CLI, `app` mediates this).
- `internal/app` (Application Layer) orchestrates the interaction between the Infrastructure (`config`, `fs`) and the Presentation (`tui`).

## Benefits
1.  **Independent of Frameworks**: The TUI library (Bubble Tea) is isolated in `internal/tui`. We could swap it for a web server or a GUI without changing the `domain` or `todo` logic.
2.  **Testable**: The `internal/domain` and `internal/todo` logic can be unit tested without any filesystem or UI.
3.  **Independent of Database**: The persistence mechanism (Markdown files) is isolated. We could switch to SQLite by changing the persistence adapter without affecting the TUI or Domain.
