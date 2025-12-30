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
	
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A cross-platform terminal-based daily journaling application.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
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

	app.Run()
}
