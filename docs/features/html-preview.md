# Feature: HTML Preview

The HTML Preview feature allows users to view their daily journal entries as beautifully formatted HTML pages in their default web browser.

## Usage

```bash
# Preview today's entry
journal preview

# Preview a specific date
journal preview 2026-01-06
```

## Functionality

1.  **Reads** the specified journal markdown file.
2.  **Parses** the content, properly handling frontmatter (hiding or formatting it).
3.  **Renders** the Markdown to HTML.
4.  **Injects** a clean, responsive CSS stylesheet (mimicking a "Reader Mode" or "Notion-like" aesthetic).
5.  **Writes** the result to a temporary HTML file.
6.  **Opens** the file using the system's default browser.

## Design

-   **Command**: `preview` (Subcommand pattern).
-   **Style**: Minimalist, centered layout, optimized for reading.
-   **Dependencies**: Requires a markdown-to-html library (e.g., `github.com/yuin/goldmark`) to be added.

## Future Enhancements

-   Live reload (server mode) while editing.
-   Export to PDF (via browser print).
-   Dark/Light mode toggle.
