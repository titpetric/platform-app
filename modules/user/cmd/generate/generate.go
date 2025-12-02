package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

var funcMap = template.FuncMap{
	"toJSON": func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	},
	"inc": func(i int) int { return i + 1 },
	"sub": func(a, b int) int { return a - b },
}

func GenerateRego(cfg *Config, outputPath string) error {
	tmpl, err := template.New("flows.rego.tmpl").Funcs(funcMap).ParseFiles("opa/flows.rego.tmpl")
	if err != nil {
		return fmt.Errorf("error in parsefiles: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return fmt.Errorf("error in execute: %w", err)
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0o644)
}
