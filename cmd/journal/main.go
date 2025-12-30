package main

import (
	"flag"
	"fmt"
	"os"

	"journal-cli/internal/app"
)

const Version = "1.0.0"

func main() {
	help := flag.Bool("help", false, "Show help message")
	version := flag.Bool("version", false, "Show version")
	todos := flag.String("todos", "", "Update todos for a date (YYYY-MM-DD). Empty = today")
	todoFlag := flag.Bool("todo", false, "Update today's todos (shorthand for --todos \"\")")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A cross-platform terminal-based daily journaling application.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nConfiguration:\n")
		fmt.Fprintf(os.Stderr, "  The application looks for a config.yaml file in:\n")
		fmt.Fprintf(os.Stderr, "  - macOS:   ~/Library/Application Support/journal-cli/config.yaml\n")
		fmt.Fprintf(os.Stderr, "  - Linux:   ~/.config/journal-cli/config.yaml\n")
		fmt.Fprintf(os.Stderr, "  - Windows: %%APPDATA%%\\journal-cli\\config.yaml\n")
		fmt.Fprintf(os.Stderr, "\n  Example config.yaml:\n")
		fmt.Fprintf(os.Stderr, "    obsidian_vault: \"/Users/username/Documents/ObsidianVault\"\n")
		fmt.Fprintf(os.Stderr, "    journal_dir: \"Journal/Daily\" # Relative to obsidian_vault\n\n")
		fmt.Fprintf(os.Stderr, "Templates:\n")
		fmt.Fprintf(os.Stderr, "  Templates are YAML files stored in the 'templates' subdirectory of the config folder.\n")
		fmt.Fprintf(os.Stderr, "  Example template:\n")
		fmt.Fprintf(os.Stderr, "    name: daily-reflection\n")
		fmt.Fprintf(os.Stderr, "    description: A simple daily reflection\n")
		fmt.Fprintf(os.Stderr, "    questions:\n")
		fmt.Fprintf(os.Stderr, "      - id: gratitude\n")
		fmt.Fprintf(os.Stderr, "        title: \"What are you grateful for?\"\n")

		fmt.Fprintf(os.Stderr, "\nTodo updater:\n")
		fmt.Fprintf(os.Stderr, "  Use --todos [YYYY-MM-DD] to run a quick CLI updater for todos (empty = today).\n")
		fmt.Fprintf(os.Stderr, "  Examples:\n")
		fmt.Fprintf(os.Stderr, "    ./journal --todos \"\"    # update today's todos\n")
		fmt.Fprintf(os.Stderr, "    ./journal --todos 2025-12-30  # update todos for that date\n")
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *version {
		fmt.Printf("journal-cli version %s\n", Version)
		return
	}

	// Run the todo updater only when explicitly requested
	if *todos != "" || *todoFlag {
		// If --todo boolean is set, pass empty string to mean today
		arg := *todos
		if *todoFlag {
			arg = ""
		}
		if err := app.UpdateTodos(arg); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating todos: %v\n", err)
			os.Exit(1)
		}
		return
	}

	app.Run()
}
