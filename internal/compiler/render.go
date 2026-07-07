package compiler

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// TemplateData holds the data passed to templates
type TemplateData struct {
	Title       string
	Description string
	Body        template.HTML
	Record      map[string]any
	Site        map[string]any
	List        []map[string]any
}

// SafeHTML helper
func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

// RenderTemplate parses the required templates and executes the specified page template
func RenderTemplate(pageName string, data TemplateData) (string, error) {
	// Parse all layouts and partials
	layouts, err := filepath.Glob("ui/templates/layouts/*.html")
	if err != nil {
		return "", err
	}
	partials, err := filepath.Glob("ui/templates/partials/*.html")
	if err != nil {
		return "", err
	}
	
	pagePath := filepath.Join("ui/templates/pages", pageName)

	var allFiles []string
	allFiles = append(allFiles, layouts...)
	allFiles = append(allFiles, partials...)
	allFiles = append(allFiles, pagePath)

	tmpl, err := template.New(filepath.Base(pagePath)).Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"hasPrefix": strings.HasPrefix,
	}).ParseFiles(allFiles...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	// Execute the "base" layout, which should embed the page content
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteStaticPage writes the generated HTML to the pb_public/cached path
func WriteStaticPage(collection, slug, html string) error {
	path := filepath.Join("pb_public", "cached", collection)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	
	filePath := filepath.Join(path, slug+".html")
	return os.WriteFile(filePath, []byte(html), 0644)
}

func DeleteStaticPage(collection, slug string) error {
	filePath := filepath.Join("pb_public", "cached", collection, slug+".html")
	return os.Remove(filePath)
}
