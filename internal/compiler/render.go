package compiler

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

var (
	tmplCacheMu sync.RWMutex
	tmplCache   = make(map[string]*template.Template)
)

// SafeHTML helper
func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

func loadTemplate(pageName string) (*template.Template, error) {
	if filepath.Base(pageName) != pageName || strings.Contains(pageName, "..") {
		return nil, fmt.Errorf("invalid template name: %s", pageName)
	}

	useCache := os.Getenv("ENV") != "development"
	if useCache {
		tmplCacheMu.RLock()
		if cached, ok := tmplCache[pageName]; ok {
			tmplCacheMu.RUnlock()
			return cached, nil
		}
		tmplCacheMu.RUnlock()
	}

	layouts, err := filepath.Glob("ui/templates/layouts/*.html")
	if err != nil {
		return nil, err
	}
	partials, err := filepath.Glob("ui/templates/partials/*.html")
	if err != nil {
		return nil, err
	}

	pagePath := filepath.Join("ui/templates/pages", pageName)

	var allFiles []string
	allFiles = append(allFiles, layouts...)
	allFiles = append(allFiles, partials...)
	allFiles = append(allFiles, pagePath)

	tmpl, err := template.New(filepath.Base(pagePath)).ParseFiles(allFiles...)
	if err != nil {
		return nil, err
	}

	if useCache {
		tmplCacheMu.Lock()
		tmplCache[pageName] = tmpl
		tmplCacheMu.Unlock()
	}

	return tmpl, nil
}

// RenderTemplate parses (and caches) the required templates and executes the specified page template
func RenderTemplate(pageName string, data TemplateData) (string, error) {
	tmpl, err := loadTemplate(pageName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	// Execute the "base" layout, which embeds the page content
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
