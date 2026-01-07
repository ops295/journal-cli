package help

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed help.yaml
var helpFS embed.FS

type Command struct {
	Name        string   `yaml:"name"`
	Args        string   `yaml:"args"`
	Description string   `yaml:"description"`
	Examples    []string `yaml:"examples"`
}

type Location struct {
	Platform string `yaml:"platform"`
	Path     string `yaml:"path"`
}

type SectionDetail struct {
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
	Locations   []Location `yaml:"locations"`
	Example     string     `yaml:"example"`
}

type ExtraSection struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Items       []string `yaml:"items"`
}

type HelpDoc struct {
	App struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
	} `yaml:"app"`
	Commands      []Command      `yaml:"commands"`
	Configuration SectionDetail  `yaml:"configuration"`
	Templates     SectionDetail  `yaml:"templates"`
	Sections      []ExtraSection `yaml:"sections"`
}

func LoadHelp() (*HelpDoc, error) {
	data, err := helpFS.ReadFile("help.yaml")
	if err != nil {
		return nil, err
	}

	var doc HelpDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (h *HelpDoc) Render(programName string) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("Usage: %s [options]\n\n", programName))
	sb.WriteString(fmt.Sprintf("%s\n\n", h.App.Description))

	// Options
	sb.WriteString("Options:\n")
	for _, cmd := range h.Commands {
		args := ""
		if cmd.Args != "" {
			args = " " + cmd.Args
		}
		// Format:   -flag args    Description
		sb.WriteString(fmt.Sprintf("  %-30s %s\n", cmd.Name+args, cmd.Description))
	}
	sb.WriteString("\n")

	// Configuration
	sb.WriteString(fmt.Sprintf("%s:\n", h.Configuration.Title))
	sb.WriteString(fmt.Sprintf("  %s:\n", h.Configuration.Description))
	for _, loc := range h.Configuration.Locations {
		sb.WriteString(fmt.Sprintf("  - %-8s %s\n", loc.Platform+":", loc.Path))
	}
	sb.WriteString("\n  Example config.yaml:\n")
	lines := strings.Split(strings.TrimSpace(h.Configuration.Example), "\n")
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("    %s\n", line))
	}
	sb.WriteString("\n")

	// Templates
	sb.WriteString(fmt.Sprintf("%s:\n", h.Templates.Title))
	sb.WriteString(fmt.Sprintf("  %s\n", h.Templates.Description))
	sb.WriteString("  Example template:\n")
	lines = strings.Split(strings.TrimSpace(h.Templates.Example), "\n")
	for _, line := range lines {
		sb.WriteString(fmt.Sprintf("    %s\n", line))
	}

	// Extra Sections (Template Management, Todo updates)
	for _, section := range h.Sections {
		sb.WriteString(fmt.Sprintf("\n%s:\n", section.Title))
		if section.Description != "" {
			sb.WriteString(fmt.Sprintf("  %s\n", section.Description))
		}
		for _, item := range section.Items {
			sb.WriteString(fmt.Sprintf("  %s\n", item))
		}

		// Find commands related to this section to show their examples
		// This is a heuristic; simpler is to just list examples from commands if they exist
		if section.Title == "Template Management" {
			sb.WriteString("  Examples:\n")
			for _, cmd := range h.Commands {
				if strings.Contains(cmd.Name, "template") && len(cmd.Examples) > 0 {
					for _, ex := range cmd.Examples {
						sb.WriteString(fmt.Sprintf("    %s\n", ex))
					}
				}
			}
		} else if section.Title == "Todo Updater" {
			sb.WriteString("  Examples:\n")
			for _, cmd := range h.Commands {
				if strings.Contains(cmd.Name, "todo") && len(cmd.Examples) > 0 {
					for _, ex := range cmd.Examples {
						sb.WriteString(fmt.Sprintf("    %s\n", ex))
					}
				}
			}
		} else if section.Title == "Updates" {
			sb.WriteString("  Examples:\n")
			for _, cmd := range h.Commands {
				if strings.Contains(cmd.Name, "update") && len(cmd.Examples) > 0 {
					for _, ex := range cmd.Examples {
						sb.WriteString(fmt.Sprintf("    %s\n", ex))
					}
				}
			}
		}
	}

	return sb.String()
}
