package main

import (
	"flag"
	"fmt"
	"os"

	"journal-cli/internal/app"
	"journal-cli/internal/help"
	"journal-cli/internal/updater"
)

const Version = "0.2.02"

func main() {
	// Check for "self-update" subcommand
	if len(os.Args) > 1 && os.Args[1] == "self-update" {
		if err := updater.Update(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
			os.Exit(1)
		}
	}

	// Check for "preview" subcommand
	if len(os.Args) > 1 && os.Args[1] == "preview" {
		dateStr := ""
		if len(os.Args) > 2 {
			dateStr = os.Args[2]
		}
		if err := app.PreviewEntry(dateStr); err != nil {
			fmt.Fprintf(os.Stderr, "Error previewing entry: %v\n", err)
			os.Exit(1)
		}
		return
	}

	helpFlag := flag.Bool("help", false, "Show help message")
	version := flag.Bool("version", false, "Show version")
	todos := flag.String("todos", "", "Update todos for a date (YYYY-MM-DD). Empty = today")
	todoFlag := flag.Bool("todo", false, "Update today's todos (shorthand for --todos \"\")")
	setTemplate := flag.String("set-template", "", "Set default template")
	listTemplates := flag.Bool("list-templates", false, "List available templates")

	// Custom usage message using YAML documentation
	flag.Usage = func() {
		helpDoc, err := help.LoadHelp()
		if err != nil {
			// Fallback if help loading fails
			fmt.Fprintf(os.Stderr, "Error loading help documentation: %v\n", err)
			fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
			flag.PrintDefaults()
			return
		}
		fmt.Fprintln(os.Stderr, helpDoc.Render(os.Args[0]))
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *version {
		fmt.Printf("journal-cli version %s\n", Version)
		return
	}

	// Handle template management commands
	if *listTemplates {
		if err := app.ListTemplates(); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing templates: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *setTemplate != "" {
		if err := app.SetDefaultTemplate(*setTemplate); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting template: %v\n", err)
			os.Exit(1)
		}
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
