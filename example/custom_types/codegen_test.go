package custom_types

import (
	"github.com/djsavvy/pggen"
	"github.com/djsavvy/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerate_Go_Example_CustomTypes(t *testing.T) {
	conn, cleanupFunc := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanupFunc()

	tmpDir := t.TempDir()
	err := pggen.Generate(
		pggen.GenerateOptions{
			ConnString: conn.Config().ConnString(),
			QueryFiles: []string{"query.sql"},
			OutputDir:  tmpDir,
			GoPackage:  "custom_types",
			Language:   pggen.LangGo,
			TypeOverrides: map[string]string{
				"text":    "github.com/djsavvy/pggen/example/custom_types/mytype.String",
				"int8":    "github.com/djsavvy/pggen/example/custom_types.CustomInt",
				"my_int":  "int",
				"_my_int": "[]int",
			},
		})
	if err != nil {
		t.Fatalf("Generate() example/custom_types: %s", err)
	}

	wantQueriesFile := "query.sql.go"
	gotQueriesFile := filepath.Join(tmpDir, "query.sql.go")
	assert.FileExists(t, gotQueriesFile, "Generate() should emit query.sql.go")
	wantQueries, err := ioutil.ReadFile(wantQueriesFile)
	if err != nil {
		t.Fatalf("read wanted query.go.sql: %s", err)
	}
	gotQueries, err := ioutil.ReadFile(gotQueriesFile)
	if err != nil {
		t.Fatalf("read generated query.go.sql: %s", err)
	}
	assert.Equalf(t, string(wantQueries), string(gotQueries),
		"Got file %s; does not match contents of %s",
		gotQueriesFile, wantQueriesFile)
}
