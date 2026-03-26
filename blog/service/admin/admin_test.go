package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		wantPage     int
		wantPageSize int
	}{
		{
			name:         "default values",
			query:        "",
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "custom page",
			query:        "page=5",
			wantPage:     5,
			wantPageSize: 10,
		},
		{
			name:         "custom pageSize",
			query:        "pageSize=25",
			wantPage:     1,
			wantPageSize: 25,
		},
		{
			name:         "both custom",
			query:        "page=3&pageSize=50",
			wantPage:     3,
			wantPageSize: 50,
		},
		{
			name:         "invalid page",
			query:        "page=invalid",
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "zero page",
			query:        "page=0",
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "negative page",
			query:        "page=-1",
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "pageSize too large",
			query:        "pageSize=200",
			wantPage:     1,
			wantPageSize: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.query, nil)
			page, pageSize := parsePagination(req)
			assert.Equal(t, tt.wantPage, page)
			assert.Equal(t, tt.wantPageSize, pageSize)
		})
	}
}
