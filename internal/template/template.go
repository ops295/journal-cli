package template

import (
	"embed"
	"os"
	"path/filepath"

	"journal-cli/internal/fs"

	"gopkg.in/yaml.v3"
)

//go:embed defaults/*.yaml
var defaultTemplatesFS embed.FS

type Question struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
}

type Template struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Questions   []Question `yaml:"questions"`
}

func LoadTemplates() ([]Template, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	templatesDir := filepath.Join(configDir, "journal-cli", "templates")
	if err := fs.EnsureDir(templatesDir); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	var templates []Template

	// If no templates found, create defaults from embedded files
	if len(files) == 0 {
		entries, err := defaultTemplatesFS.ReadDir("defaults")
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			data, err := defaultTemplatesFS.ReadFile("defaults/" + entry.Name())
			if err != nil {
				continue
			}

			// Write to config dir
			if err := fs.WriteFile(filepath.Join(templatesDir, entry.Name()), data); err != nil {
				continue
			}

			var t Template
			if err := yaml.Unmarshal(data, &t); err != nil {
				continue
			}
			templates = append(templates, t)
		}
		return templates, nil
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			data, err := os.ReadFile(filepath.Join(templatesDir, file.Name()))
			if err != nil {
				continue
			}
			var t Template
			if err := yaml.Unmarshal(data, &t); err != nil {
				continue
			}
			templates = append(templates, t)
		}
	}

	return templates, nil
}
