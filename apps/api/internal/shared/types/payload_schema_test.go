package types

import (
	"encoding/json"
	"testing"
)

const invoiceSchema = `{
	"type": "object",
	"properties": {
		"invoice_id": {"type": "string"},
		"amount_cents": {"type": "integer", "minimum": 0}
	},
	"required": ["invoice_id", "amount_cents"]
}`

func TestNewPayloadSchema(t *testing.T) {
	t.Run("accepts a valid JSON Schema", func(t *testing.T) {
		if _, err := NewPayloadSchema(json.RawMessage(invoiceSchema)); err != nil {
			t.Fatalf("expected schema to be accepted, got: %v", err)
		}
	})

	t.Run("rejects malformed JSON", func(t *testing.T) {
		if _, err := NewPayloadSchema(json.RawMessage(`{not json`)); err == nil {
			t.Fatal("expected error for malformed JSON")
		}
	})

	t.Run("rejects an invalid JSON Schema", func(t *testing.T) {
		if _, err := NewPayloadSchema(json.RawMessage(`{"type": "not-a-type"}`)); err == nil {
			t.Fatal("expected error for invalid JSON Schema")
		}
	})
}

func TestPayloadSchemaValidate(t *testing.T) {
	schema, err := NewPayloadSchema(json.RawMessage(invoiceSchema))
	if err != nil {
		t.Fatalf("schema should compile: %v", err)
	}

	t.Run("accepts a conforming payload", func(t *testing.T) {
		payload := []byte(`{"invoice_id": "inv_123", "amount_cents": 4200}`)
		if err := schema.Validate(payload); err != nil {
			t.Fatalf("expected payload to validate, got: %v", err)
		}
	})

	t.Run("rejects a payload missing required fields", func(t *testing.T) {
		if err := schema.Validate([]byte(`{"invoice_id": "inv_123"}`)); err == nil {
			t.Fatal("expected error for missing amount_cents")
		}
	})

	t.Run("rejects a payload with wrong types", func(t *testing.T) {
		payload := []byte(`{"invoice_id": "inv_123", "amount_cents": "a lot"}`)
		if err := schema.Validate(payload); err == nil {
			t.Fatal("expected error for string amount_cents")
		}
	})

	t.Run("rejects malformed JSON payloads", func(t *testing.T) {
		if err := schema.Validate([]byte(`{broken`)); err == nil {
			t.Fatal("expected error for malformed payload")
		}
	})
}
