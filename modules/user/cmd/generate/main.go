package main

import (
	"log"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	if err := ValidateFlowsYAML("opa/flows.yml", "opa/flows.jsonschema"); err != nil {
		return err
	}

	cfg, err := LoadConfig("opa/flows.yml")
	if err != nil {
		return err
	}

	if err := GenerateRego(cfg, "opa/flows.rego"); err != nil {
		return err
	}
	if err := GeneratePlantUML(cfg, "opa/flows.puml"); err != nil {
		return err
	}
	if err := GenerateRouteMap(cfg, "opa/routes.json"); err != nil {
		return err
	}

	log.Println("Generated: opa/flows.rego")
	return nil
}
