package blog

// RenderData prepares form data for vuego template rendering
type RenderData struct {
	Title       string
	Description string
	Fields      []Field
	Method      string
	Action      string
	SubmitLabel string
	CancelURL   string
	CancelLabel string
	Navigation  *Navigation
}

// NewRenderData creates render data from a form and navigation
func NewRenderData(form *Form, nav *Navigation) *RenderData {
	return &RenderData{
		Title:       form.Title,
		Fields:      form.Fields,
		Method:      "POST",
		Action:      "#",
		SubmitLabel: "Save",
		CancelLabel: "Cancel",
		Navigation:  nav,
	}
}

// SetAction sets the form action URL
func (rd *RenderData) SetAction(action string) *RenderData {
	rd.Action = action
	return rd
}

// SetMethod sets the HTTP method (GET, POST, etc)
func (rd *RenderData) SetMethod(method string) *RenderData {
	rd.Method = method
	return rd
}

// SetSubmitLabel sets the submit button label
func (rd *RenderData) SetSubmitLabel(label string) *RenderData {
	rd.SubmitLabel = label
	return rd
}

// SetCancelURL sets the cancel button URL
func (rd *RenderData) SetCancelURL(url string) *RenderData {
	rd.CancelURL = url
	return rd
}

// SetCancelLabel sets the cancel label
func (rd *RenderData) SetCancelLabel(label string) *RenderData {
	rd.CancelLabel = label
	return rd
}

// VisibleFields returns only non-hidden fields
func (rd *RenderData) VisibleFields() []Field {
	var visible []Field
	hiddenTypes := map[string]bool{
		"hidden": true, "key": true, "parent_key": true, "explicit_parent_key": true,
	}
	for _, f := range rd.Fields {
		if !hiddenTypes[f.Type] {
			visible = append(visible, f)
		}
	}
	return visible
}

// HiddenFields returns only hidden/key fields
func (rd *RenderData) HiddenFields() []Field {
	var hidden []Field
	hiddenTypes := map[string]bool{
		"hidden": true, "key": true, "parent_key": true, "explicit_parent_key": true,
	}
	for _, f := range rd.Fields {
		if hiddenTypes[f.Type] {
			hidden = append(hidden, f)
		}
	}
	return hidden
}

// FieldsBySize groups fields by size
func (rd *RenderData) FieldsBySize(size string) []Field {
	var result []Field
	for _, f := range rd.Fields {
		if f.Size == size {
			result = append(result, f)
		}
	}
	return result
}

// FieldsWithType returns fields of a specific type
func (rd *RenderData) FieldsWithType(fieldType string) []Field {
	var result []Field
	for _, f := range rd.Fields {
		if f.Type == fieldType {
			result = append(result, f)
		}
	}
	return result
}
