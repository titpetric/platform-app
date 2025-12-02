package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"text/template"
)

func GeneratePlantUML(cfg *Config, outputPath string) error {
	tmpl, err := template.New("flows.puml.tmpl").Funcs(funcMap).ParseFiles("opa/flows.puml.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return err
	}

	// Write .puml file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
		return err
	}

	log.Println("Generated PlantUML file:", outputPath)

	// Run plantuml CLI to generate SVG
	cmd := exec.Command("plantuml", "-tsvg", outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	log.Println("Generated SVG from PlantUML:", outputPath[:len(outputPath)-5]+".svg") // remove .puml
	return nil
}
