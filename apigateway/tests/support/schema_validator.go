package support

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidator valida datos contra esquemas JSON-Schema
type SchemaValidator struct {
	schemas map[string]*gojsonschema.Schema
}

// NewSchemaValidator crea un nuevo validador
func NewSchemaValidator(schemasPath string) (*SchemaValidator, error) {
	sv := &SchemaValidator{
		schemas: make(map[string]*gojsonschema.Schema),
	}

	// Leer archivo de esquemas
	schemasData, err := os.ReadFile(schemasPath)
	if err != nil {
		return nil, fmt.Errorf("error reading schemas file: %w", err)
	}

	// Parsear esquemas
	var schemasMap map[string]interface{}
	if err := json.Unmarshal(schemasData, &schemasMap); err != nil {
		return nil, fmt.Errorf("error parsing schemas: %w", err)
	}

	// Compilar cada esquema
	for schemaName, schemaData := range schemasMap {
		schemaJSON, err := json.Marshal(schemaData)
		if err != nil {
			return nil, fmt.Errorf("error marshaling schema %s: %w", schemaName, err)
		}

		schema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(string(schemaJSON)))
		if err != nil {
			return nil, fmt.Errorf("error compiling schema %s: %w", schemaName, err)
		}

		sv.schemas[schemaName] = schema
	}

	return sv, nil
}

// ValidationResult encapsula el resultado de una validación
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// Validate valida datos contra un esquema
func (sv *SchemaValidator) Validate(data interface{}, schemaName string) *ValidationResult {
	schema, ok := sv.schemas[schemaName]
	if !ok {
		return &ValidationResult{
			IsValid: false,
			Errors:  []string{fmt.Sprintf("Schema '%s' not found", schemaName)},
		}
	}

	dataLoader := gojsonschema.NewGoLoader(data)
	result, err := schema.Validate(dataLoader)

	if err != nil {
		return &ValidationResult{
			IsValid: false,
			Errors:  []string{err.Error()},
		}
	}

	if !result.Valid() {
		var errors []string
		for _, err := range result.Errors() {
			errors = append(errors, err.String())
		}
		return &ValidationResult{
			IsValid: false,
			Errors:  errors,
		}
	}

	return &ValidationResult{
		IsValid: true,
		Errors:  []string{},
	}
}

// ValidateJSON valida JSON string contra un esquema
func (sv *SchemaValidator) ValidateJSON(jsonStr string, schemaName string) *ValidationResult {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return &ValidationResult{
			IsValid: false,
			Errors:  []string{fmt.Sprintf("Invalid JSON: %v", err)},
		}
	}

	return sv.Validate(data, schemaName)
}
