package types

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// PayloadSchema is a JSON Schema document that event_type payloads are validated
// against at ingest time. Construct via NewPayloadSchema so an invalid
// schema can never exist.
type PayloadSchema struct {
	raw      json.RawMessage
	compiled *jsonschema.Schema
}

func NewPayloadSchema(raw json.RawMessage) (PayloadSchema, error) {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	if err != nil {
		return PayloadSchema{}, fmt.Errorf("payload schema is not valid JSON: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", doc); err != nil {
		return PayloadSchema{}, fmt.Errorf("payload schema is not a valid JSON Schema: %w", err)
	}

	compiled, err := compiler.Compile("schema.json")
	if err != nil {
		return PayloadSchema{}, fmt.Errorf("payload schema is not a valid JSON Schema: %w", err)
	}

	return PayloadSchema{raw: raw, compiled: compiled}, nil
}

// Validate reports whether payload conforms to the schema.
func (s PayloadSchema) Validate(payload []byte) error {
	instance, err := jsonschema.UnmarshalJSON(bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("payload is not valid JSON: %w", err)
	}

	return s.compiled.Validate(instance)
}

// JSON returns the original schema document, e.g. for persisting to the
// payload_schema jsonb column.
func (s PayloadSchema) JSON() json.RawMessage {
	return s.raw
}
