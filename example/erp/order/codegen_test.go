package order

import (
	"github.com/djsavvy/pggen"
	"github.com/djsavvy/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerate_Go_Example_Order(t *testing.T) {
	conn, cleanupFunc := pgtest.NewPostgresSchema(t, []string{
		"../01_schema.sql",
		"../02_schema.sql",
	})
	defer cleanupFunc()

	tmpDir := t.TempDir()
	err := pggen.Generate(
		pggen.GenerateOptions{
			ConnString: conn.Config().ConnString(),
			QueryFiles: []string{
				"customer.sql",
				"price.sql",
			},
			OutputDir:     tmpDir,
			GoPackage:     "order",
			Language:      pggen.LangGo,
			Acronyms:      map[string]string{"mrr": "MRR"},
			TypeOverrides: map[string]string{"tenant_id": "int"},
		})
	if err != nil {
		t.Fatalf("Generate() example/erp/order: %s", err)
	}

	for _, file := range []string{"customer.sql.go", "price.sql.go"} {
		wantQueries, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("read wanted file %s: %s", file, err)
		}

		gotFile := filepath.Join(tmpDir, file)
		assert.FileExists(t, gotFile, "Generate() should emit "+file)
		gotQueries, err := ioutil.ReadFile(gotFile)
		if err != nil {
			t.Fatalf("read generated %s: %s", file, err)
		}
		assert.Equalf(t, string(wantQueries), string(gotQueries),
			"Got file %s; does not match contents of file %s",
			gotFile, file)
	}
}
