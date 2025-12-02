package main

import (
	"encoding/json"
	"log"
	"os"
)

func GenerateRouteMap(cfg *Config, outputPath string) error {
	routeMap := make(map[string]string)

	for flowName, flow := range cfg.Flows {
		for _, step := range flow.Steps {
			routeKey := flowName + "." + step.Name
			routeMap[routeKey] = step.Link
		}
	}

	data, err := json.MarshalIndent(routeMap, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return err
	}

	log.Println("Generated route map:", outputPath)
	return nil
}
