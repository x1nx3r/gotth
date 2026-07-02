package docs

import (
	"fmt"
	"net/http"
	"strings"
	"gotth/app/assets"
)

func PageHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	slug := strings.TrimPrefix(path, "/docs/")
	if slug == "/docs" || slug == "" || slug == "docs" {
		slug = "getting-started"
	}

	// Read markdown file from virtual embedded filesystem
	content, err := assets.Docs.ReadFile("docs/" + slug + ".md")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Simple Custom Markdown-to-HTML parser (Zero dependency)
	html := parseMarkdown(string(content))
	
	// Convert hyphenated slug to friendly title
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	title := strings.Join(words, " ")

	Page(title, html, "/docs/"+slug).Render(r.Context(), w)
}

func parseMarkdown(md string) string {
	var html strings.Builder
	lines := strings.Split(md, "\n")
	inCodeBlock := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Handle Code Blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				html.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				html.WriteString("<pre><code>")
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			html.WriteString(lines[i] + "\n")
			continue
		}

		// Handle Headers
		if strings.HasPrefix(line, "# ") {
			html.WriteString(fmt.Sprintf("<h1>%s</h1>\n", line[2:]))
		} else if strings.HasPrefix(line, "## ") {
			html.WriteString(fmt.Sprintf("<h2>%s</h2>\n", line[3:]))
		} else {
			// Handle Paragraphs and inline backticks `code`
			parsedLine := line
			for {
				start := strings.Index(parsedLine, "`")
				if start == -1 {
					break
				}
				end := strings.Index(parsedLine[start+1:], "`")
				if end == -1 {
					break
				}
				endIdx := start + 1 + end
				codeContent := parsedLine[start+1 : endIdx]
				parsedLine = parsedLine[:start] + "<code>" + codeContent + "</code>" + parsedLine[endIdx+1:]
			}
			html.WriteString(fmt.Sprintf("<p>%s</p>\n", parsedLine))
		}
	}
	return html.String()
}
