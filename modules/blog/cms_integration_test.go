package blog

import (
	"testing"
)

func TestCMSFormTypes(t *testing.T) {
	// Test Form type creation
	form := &Form{
		Name:   "test_form",
		Title:  "Test Form",
		Fields: []Field{},
	}

	if form.Name != "test_form" {
		t.Errorf("expected form name 'test_form', got '%s'", form.Name)
	}

	// Test Field type creation
	field := Field{
		Name:     "test_field",
		Title:    "Test Field",
		Type:     "text",
		Required: true,
		Value:    "test_value",
	}

	if field.Type != "text" {
		t.Errorf("expected field type 'text', got '%s'", field.Type)
	}

	if field.Value != "test_value" {
		t.Errorf("expected field value 'test_value', got '%v'", field.Value)
	}

	// Test Option type creation
	option := Option{
		Value: "opt1",
		Label: "Option 1",
	}

	if option.Value != "opt1" {
		t.Errorf("expected option value 'opt1', got '%s'", option.Value)
	}

	// Test ParseFieldValue with nil
	opts := ParseFieldValue(nil)
	if opts != nil {
		t.Error("expected nil for ParseFieldValue(nil)")
	}
}

func TestCMSRenderData(t *testing.T) {
	form := &Form{
		Name:  "test_form",
		Title: "Test Form",
		Fields: []Field{
			{
				Name: "title",
				Type: "text",
			},
			{
				Name: "content",
				Type: "html",
			},
			{
				Name: "hidden_id",
				Type: "hidden",
			},
			{
				Name: "db_key",
				Type: "key",
			},
		},
	}

	rd := NewRenderData(form, nil)

	if rd.Title != "Test Form" {
		t.Errorf("expected title 'Test Form', got '%s'", rd.Title)
	}

	visibleFields := rd.VisibleFields()
	if len(visibleFields) != 2 {
		t.Errorf("expected 2 visible fields, got %d", len(visibleFields))
	}

	hiddenFields := rd.HiddenFields()
	if len(hiddenFields) != 2 {
		t.Errorf("expected 2 hidden fields, got %d", len(hiddenFields))
	}

	// Test fluent API
	rd.SetAction("/save").SetSubmitLabel("Publish")
	if rd.Action != "/save" {
		t.Errorf("expected action '/save', got '%s'", rd.Action)
	}
	if rd.SubmitLabel != "Publish" {
		t.Errorf("expected submit label 'Publish', got '%s'", rd.SubmitLabel)
	}
}
