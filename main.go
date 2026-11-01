// Static site generator for the resume site.
//
// Reads data/resume.json, renders templates/index.html and copies static/
// into public/. Run with: go run .
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type SkillGroup struct {
	Label string   `json:"label"`
	Items []string `json:"items"`
}

type Role struct {
	Title   string   `json:"title"`
	Period  string   `json:"period"`
	Project string   `json:"project"`
	Bullets []string `json:"bullets"`
}

type Company struct {
	Name  string `json:"name"`
	Roles []Role `json:"roles"`
}

type Education struct {
	Degree string `json:"degree"`
	School string `json:"school"`
	Period string `json:"period"`
}

type Resume struct {
	Name        string       `json:"name"`
	Title       string       `json:"title"`
	Location    string       `json:"location"`
	Tagline     string       `json:"tagline"`
	Description string       `json:"description"`
	Summary     []string     `json:"summary"`
	Email       string       `json:"email"`
	Links       []Link       `json:"links"`
	Skills      []SkillGroup `json:"skills"`
	Experience  []Company    `json:"experience"`
	Education   []Education  `json:"education"`
}

const (
	dataFile    = "data/resume.json"
	templateDir = "templates"
	staticDir   = "static"
	outDir      = "public"
)

func main() {
	resume, err := loadResume(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	warnPlaceholders(resume)

	if err := os.RemoveAll(outDir); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	if err := copyDir(staticDir, outDir); err != nil {
		log.Fatal(err)
	}

	if err := render(resume); err != nil {
		log.Fatal(err)
	}

	// Tell GitHub Pages not to run Jekyll over the output.
	if err := os.WriteFile(filepath.Join(outDir, ".nojekyll"), nil, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("built %s/index.html\n", outDir)
}

func loadResume(path string) (*Resume, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var r Resume
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &r, nil
}

func render(r *Resume) error {
	funcs := template.FuncMap{
		"join": strings.Join,
	}
	tmpl, err := template.New("index.html").
		Funcs(funcs).
		ParseFiles(filepath.Join(templateDir, "index.html"))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, r)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// warnPlaceholders flags profile links that still need real values.
func warnPlaceholders(r *Resume) {
	for _, l := range r.Links {
		if strings.Contains(l.URL, "your-") {
			log.Printf("warning: %q link is still a placeholder (%s) - update data/resume.json", l.Label, l.URL)
		}
	}
}
