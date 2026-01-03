package schema_test

import (
	"os"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/titpetric/platform-app/modules/cms/schema"
)

//go:embed schema.yaml
var entitySchema []byte

func TestLoad(t *testing.T) {
	data, err := schema.LoadSchema(entitySchema)
	assert.NoError(t, err)

	tableNames := data.GetTableNames()
	for _, name := range tableNames {
		table := data.GetTable(name)
		t.Logf("Table %s: %s", table.Name, table.Comment)
	}

	data.Walk(func(table *schema.Table) {
		t.Logf("Table %s: %s", table.Name, table.Comment)
	})

	if testing.Verbose() && false {
		assert.NoError(t, yaml.NewEncoder(os.Stderr).Encode(data))
	}
}
