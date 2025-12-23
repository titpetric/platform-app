package cms

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFormFromYAML(t *testing.T) {
	// Load example article form
	data, err := os.ReadFile("testdata/article.yaml")
	assert.NoError(t, err, "failed to read example form")

	form, err := LoadFormFromYAML("article", data)
	assert.NoError(t, err, "failed to load form")
	assert.NotNil(t, form, "form should not be nil")

	assert.Equal(t, "article", form.Name, "form name should be 'article'")
	assert.Greater(t, len(form.Fields), 0, "form should have fields")

	// Check for title field
	titleField := findField(form.Fields, "title")
	assert.NotNil(t, titleField, "expected title field")
	assert.Equal(t, "string", titleField.Type, "title field type should be 'string'")

	// Check for status enum field
	statusField := findField(form.Fields, "status")
	assert.NotNil(t, statusField, "expected status field")
	assert.Greater(t, len(statusField.Options), 0, "status field should have parsed options")
	assert.Equal(t, "draft", statusField.Options[0].Value, "first option value should be 'draft'")

	// Check for hidden fields
	idField := findField(form.Fields, "id")
	assert.NotNil(t, idField, "expected id field")
	assert.Equal(t, "key", idField.Type, "id field type should be 'key'")
}

func TestRenderData(t *testing.T) {
	data, err := os.ReadFile("testdata/article.yaml")
	assert.NoError(t, err, "failed to read example form")

	form, err := LoadFormFromYAML("article", data)
	assert.NoError(t, err, "failed to load form")

	rd := NewRenderData(form, nil)
	rd.SetAction("/articles/save").SetSubmitLabel("Publish")

	assert.Equal(t, "/articles/save", rd.Action, "action should be '/articles/save'")
	assert.Equal(t, "Publish", rd.SubmitLabel, "submit label should be 'Publish'")

	visibleFields := rd.VisibleFields()
	hiddenFields := rd.HiddenFields()

	assert.Greater(t, len(visibleFields), 0, "should have visible fields")
	assert.Greater(t, len(hiddenFields), 0, "should have hidden fields")
	assert.Equal(t, len(form.Fields), len(visibleFields)+len(hiddenFields),
		"visible and hidden fields should equal total fields")
}

func TestLoadFormParseEnumValue(t *testing.T) {
	data, err := os.ReadFile("testdata/article.yaml")
	assert.NoError(t, err, "failed to read example form")

	form, err := LoadFormFromYAML("article", data)
	assert.NoError(t, err, "failed to load form")

	statusField := findField(form.Fields, "status")
	assert.NotNil(t, statusField, "expected status field")

	expectedOptions := []Option{
		{Value: "draft", Label: "Draft"},
		{Value: "published", Label: "Published"},
		{Value: "archived", Label: "Archived"},
	}

	assert.Equal(t, len(expectedOptions), len(statusField.Options), "should have 3 options")
	for i, expected := range expectedOptions {
		assert.Equal(t, expected.Value, statusField.Options[i].Value, "option %d value mismatch", i)
		assert.Equal(t, expected.Label, statusField.Options[i].Label, "option %d label mismatch", i)
	}
}

func TestFieldGrouping(t *testing.T) {
	data, err := os.ReadFile("testdata/article.yaml")
	assert.NoError(t, err, "failed to read example form")

	form, err := LoadFormFromYAML("article", data)
	assert.NoError(t, err, "failed to load form")

	rd := NewRenderData(form, nil)

	visibleFields := rd.VisibleFields()
	hiddenFields := rd.HiddenFields()
	fullSizeFields := rd.FieldsBySize("full")

	assert.Equal(t, len(form.Fields), len(visibleFields)+len(hiddenFields),
		"all fields should be accounted for")
	assert.Greater(t, len(fullSizeFields), 0, "should have full-size fields")
}

func TestRenderDataFluentAPI(t *testing.T) {
	form := &Form{
		Name:   "test",
		Title:  "Test Form",
		Fields: []Field{},
	}

	rd := NewRenderData(form, nil)
	result := rd.SetAction("/save").
		SetMethod("PUT").
		SetSubmitLabel("Update").
		SetCancelURL("/cancel").
		SetCancelLabel("Go Back")

	assert.Equal(t, rd, result, "fluent API should return receiver")
	assert.Equal(t, "/save", rd.Action)
	assert.Equal(t, "PUT", rd.Method)
	assert.Equal(t, "Update", rd.SubmitLabel)
	assert.Equal(t, "/cancel", rd.CancelURL)
	assert.Equal(t, "Go Back", rd.CancelLabel)
}

func findField(fields []Field, name string) *Field {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
}
