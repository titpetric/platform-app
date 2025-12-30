package model

// FormField represents a form field for rendering
type FormField struct {
	Name        string        `json:"name"`
	Label       string        `json:"label"`
	Type        string        `json:"type"` // text, number, select, textarea, checkbox, datetime-local, email, url, tel, password
	Value       interface{}   `json:"value,omitempty"`
	Placeholder string        `json:"placeholder,omitempty"`
	Required    bool          `json:"required"`
	Readonly    bool          `json:"readonly"`
	Disabled    bool          `json:"disabled"`
	Options     []FieldOption `json:"options,omitempty"`
	Rows        int           `json:"rows,omitempty"`
	Cols        int           `json:"cols,omitempty"`
	Help        string        `json:"help,omitempty"`
	Error       string        `json:"error,omitempty"`
	Size        string        `json:"size,omitempty"` // full, half
	Attributes  map[string]string `json:"attributes,omitempty"`
}

// FieldOption represents a select/radio option
type FieldOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// FormData represents a complete form for editing/creating a record
type FormData struct {
	Title      string      `json:"title"`
	TableName  string      `json:"tableName"`
	RecordID   interface{} `json:"recordId,omitempty"`
	Fields     []FormField `json:"fields"`
	SubmitURL  string      `json:"submitUrl"`
	SubmitText string      `json:"submitText"`
	CancelURL  string      `json:"cancelUrl,omitempty"`
	Method     string      `json:"method"` // POST or PUT
}

// ListData represents a paginated list with filters
type ListData struct {
	Title       string        `json:"title"`
	TableName   string        `json:"tableName"`
	Columns     []ListColumn  `json:"columns"`
	Rows        []map[string]interface{} `json:"rows"`
	CreateURL   string        `json:"createUrl"`
	EditURLBase string        `json:"editUrlBase"`
	DeleteURLBase string       `json:"deleteUrlBase"`
	Pagination  PaginationInfo `json:"pagination"`
	Filters     []FilterField `json:"filters,omitempty"`
	Search      string        `json:"search,omitempty"`
}

// ListColumn represents a column to display in the list
type ListColumn struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Type    string `json:"type"` // text, number, enum, timestamp, boolean
	Width   string `json:"width,omitempty"` // e.g., "30%", "100px"
	Sortable bool  `json:"sortable"`
}

// FilterField represents a filter option
type FilterField struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Type    string `json:"type"` // text, select, date, daterange
	Options []FieldOption `json:"options,omitempty"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
	HasPrev    bool `json:"hasPrev"`
	HasNext    bool `json:"hasNext"`
	PrevURL    string `json:"prevUrl,omitempty"`
	NextURL    string `json:"nextUrl,omitempty"`
}

// BuildFormFromTable creates a FormData from a schema table
func BuildFormFromTable(table *Table, records map[string]interface{}, isCreate bool) *FormData {
	form := &FormData{
		TableName: table.Name,
		Fields:    make([]FormField, 0),
	}

	if isCreate {
		form.Title = "Create " + table.Comment
		form.SubmitText = "Create"
		form.Method = "POST"
	} else {
		form.Title = "Edit " + table.Comment
		form.SubmitText = "Update"
		form.Method = "PUT"
		if id, ok := records["id"]; ok {
			form.RecordID = id
		}
	}

	// Build fields from columns, excluding primary keys and system columns
	for _, col := range table.Columns {
		if col.IsPrimaryKey() && !isCreate {
			continue // Don't edit primary keys
		}

		field := FormField{
			Name:        col.Name,
			Label:       col.GetLabel(),
			Type:        col.GetFieldType(),
			Required:    !col.Nullable && !col.IsPrimaryKey(),
			Help:        col.Comment,
			Value:       records[col.Name],
			Attributes:  make(map[string]string),
		}

		// Special handling for enum columns
		if col.IsEnum() {
			field.Options = make([]FieldOption, len(col.EnumValues))
			for i, val := range col.EnumValues {
				field.Options[i] = FieldOption{
					Value: val,
					Label: val,
				}
			}
		}

		// Special handling for textarea columns (long text)
		if col.DataType == "text" && !col.IsEnum() {
			field.Type = "textarea"
			field.Rows = 4
		}

		// Skip system columns for editing
		if col.Name == "created_at" || col.Name == "updated_at" {
			if !isCreate {
				field.Readonly = true
			}
		}

		form.Fields = append(form.Fields, field)
	}

	return form
}

// BuildListFromTable creates ListData from a schema table
func BuildListFromTable(table *Table, rows []map[string]interface{}, page, pageSize, total int) *ListData {
	list := &ListData{
		Title:     table.Comment + " List",
		TableName: table.Name,
		Columns:   make([]ListColumn, 0),
		Rows:      rows,
		CreateURL: "/cms/" + table.Name + "/create",
		EditURLBase: "/cms/" + table.Name + "/edit/",
		DeleteURLBase: "/cms/" + table.Name + "/delete/",
	}

	// Build columns from schema (show first 5 columns by default)
	maxCols := 5
	if len(table.Columns) < maxCols {
		maxCols = len(table.Columns)
	}

	for i := 0; i < maxCols; i++ {
		col := table.Columns[i]
		list.Columns = append(list.Columns, ListColumn{
			Name:     col.Name,
			Label:    col.GetLabel(),
			Type:     col.DataType,
			Sortable: true,
		})
	}

	// Calculate pagination
	totalPages := (total + pageSize - 1) / pageSize
	list.Pagination = PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}

	if page > 1 {
		list.Pagination.PrevURL = buildListURL(table.Name, page-1, pageSize)
	}
	if page < totalPages {
		list.Pagination.NextURL = buildListURL(table.Name, page+1, pageSize)
	}

	return list
}

func buildListURL(table string, page, pageSize int) string {
	return "/cms/" + table + "?page=" + string(rune(page)) + "&pageSize=" + string(rune(pageSize))
}
