package app

import (
	"fmt"

	"journal-cli/internal/config"
	"journal-cli/internal/template"
)

// SetDefaultTemplate sets the default template in config
func SetDefaultTemplate(templateName string) error {
	// Validate template exists
	exists, err := template.TemplateExists(templateName)
	if err != nil {
		return fmt.Errorf("failed to check template: %w", err)
	}
	if !exists {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Load current config
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Update config with default template
	cfg.DefaultTemplate = templateName

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Default template set to: %s\n", templateName)
	return nil
}

// ListTemplates displays all available templates
func ListTemplates() error {
	templates, err := template.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	// Load config to check for default template
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Available templates:")
	fmt.Println()
	for _, t := range templates {
		marker := " "
		if t.Name == cfg.DefaultTemplate {
			marker = "✓"
		}
		fmt.Printf("  %s %s\n", marker, t.Name)
		if t.Description != "" {
			fmt.Printf("    %s\n", t.Description)
		}
		fmt.Println()
	}

	if cfg.DefaultTemplate != "" {
		fmt.Printf("Default template: %s\n", cfg.DefaultTemplate)
	} else {
		fmt.Println("No default template set. Use --set-template <name> to set one.")
	}

	return nil
}

// GetDefaultTemplate returns the configured default template
func GetDefaultTemplate() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.DefaultTemplate, nil
}
