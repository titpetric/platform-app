package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v3"
)

func ValidateFlowsYAML(yamlPath, schemaPath string) error {
	// Read YAML
	dataYAML, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	var data interface{}
	if err := yaml.Unmarshal(dataYAML, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling to json: %w", err)
	}

	// Read JSON Schema
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	rs := &jsonschema.Schema{}
	if err := rs.UnmarshalJSON(schemaBytes); err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	// Validate
	errs, err := rs.ValidateBytes(context.Background(), dataJSON)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("YAML validation errors: %v", errs)
	}

	return nil
}
