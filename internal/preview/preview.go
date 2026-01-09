package preview

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// CSS for the preview. Simple, clean, centered, responsive.
const previewCSS = `
<style>
	body {
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
		line-height: 1.6;
		color: #333;
		max-width: 800px;
		margin: 0 auto;
		padding: 20px;
		background-color: #f9f9f9;
	}
	article {
		background: white;
		padding: 40px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0,0,0,0.1);
	}
	h1, h2, h3 { color: #111; }
	h1 { border-bottom: 2px solid #eaeaea; padding-bottom: 0.3em; }
	blockquote { border-left: 4px solid #ddd; padding-left: 1em; color: #666; }
	code { background: #eee; padding: 2px 4px; border-radius: 4px; }
	pre { background: #f4f4f4; padding: 15px; overflow-x: auto; border-radius: 4px; }
	pre code { background: none; padding: 0; }
	ul, ol { padding-left: 2em; }
	input[type="checkbox"] { margin-right: 0.5em; }
	.meta { color: #888; font-size: 0.9em; margin-bottom: 20px; }
	@media (prefers-color-scheme: dark) {
		body { background-color: #1a1a1a; color: #ddd; }
		article { background: #2d2d2d; color: #ddd; }
		h1, h2, h3 { color: #fff; border-color: #444; }
		code { background: #444; }
		pre { background: #333; }
		blockquote { border-color: #444; color: #aaa; }
	}
</style>
`

// Render converts markdown content to a full HTML page
func Render(markdownContent []byte, title string) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(), // Allow inline HTML if needed
		),
	)

	var buf bytes.Buffer
	buf.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	buf.WriteString("<meta charset=\"UTF-8\">\n")
	buf.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	buf.WriteString("<title>" + title + "</title>\n")
	buf.WriteString(previewCSS)
	buf.WriteString("</head>\n<body>\n<article>\n")

	// Convert Markdown
	if err := md.Convert(markdownContent, &buf); err != nil {
		return nil, err
	}

	buf.WriteString("\n</article>\n</body>\n</html>")
	return buf.Bytes(), nil
}

// OpenInBrowser saves content to a temp file and opens it
func OpenInBrowser(content []byte) error {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("journal-preview-%d.html", time.Now().Unix()))
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", tmpFile)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", tmpFile)
	default: // linux, etc
		cmd = exec.Command("xdg-open", tmpFile)
	}

	return cmd.Start()
}
