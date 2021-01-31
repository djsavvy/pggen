package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/jschaf/pggen"
	"github.com/jschaf/pggen/internal/flags"
	"github.com/peterbourgon/ff/v3/ffcli"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const flagHelp = `
pggen generates type-safe code from files containing Postgres queries by running
the queries on Postgres to get type information.

EXAMPLES
  # Generate code for a single query file using an existing postgres database.
  pggen gen go --query-glob author/queries.sql --postgres-connection "user=postgres port=5555 dbname=pggen"

  # Generate code using Docker to create the postgres database with a schema 
  # file. --schema-glob arg implies using Dockerized postgres.
  pggen gen go --schema-glob author/schema.sql --query-glob author/queries.sql

  # Generate code for all queries underneath a directory. Glob should be quoted
  # to prevent shell expansion.
  pggen gen go --schema-glob author/schema.sql --query-glob 'author/**/*.sql'
`

func run() error {
	genCmd := newGenCmd()
	rootFlagSet := flag.NewFlagSet("root", flag.ExitOnError)
	rootCmd := &ffcli.Command{
		ShortUsage: "pggen <subcommand> [options...]",
		LongHelp:   flagHelp[1 : len(flagHelp)-1], // remove lead/trail newlines
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			genCmd,
		},
	}
	rootCmd.Exec = func(ctx context.Context, args []string) error {
		fmt.Println(ffcli.DefaultUsageFunc(rootCmd))
		os.Exit(1)
		return nil
	}
	if err := rootCmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		return err
	}
	return nil
}

func newGenCmd() *ffcli.Command {
	fset := flag.NewFlagSet("go", flag.ExitOnError)
	outputDir := fset.String("output-dir", "",
		"where to write generated code; defaults to same directory as query files")
	postgresConn := fset.String("postgres-connection", "",
		`connection string to a postgres database, like: `+
			`"user=postgres host=localhost dbname=pggen"`)
	queryGlobs := flags.Strings(fset, "query-glob", nil,
		"generate code for all SQL files that match glob, like 'queries/**/*.sql'")
	schemaGlobs := flags.Strings(fset, "schema-glob", nil,
		"create schema in Dockerized Postgres from all sql, sql.gz, or shell "+
			"scripts (*.sh) that match a glob, like 'migrations/*.sql'")
	goSubCmd := &ffcli.Command{
		Name:       "go",
		ShortUsage: "pggen gen go [options...]",
		ShortHelp:  "generates go code for Postgres query files",
		FlagSet:    fset,
		Exec: func(ctx context.Context, args []string) error {
			// Preconditions.
			if len(*queryGlobs) == 0 {
				return fmt.Errorf("pggen gen go: at least one file in --query-glob must match")
			}
			if *schemaGlobs != nil && *postgresConn != "" {
				return fmt.Errorf("cannot use both --schema-glob and --postgres-connection together\n" +
					"    use --schema-glob to run dockerized postgres automatically\n" +
					"    use --postgres-connection to connect to an existing database")
			}

			// Get absolute paths for all files matching each --query-glob.
			queries := make([]string, 0, len(*queryGlobs)*4)
			for _, glob := range *queryGlobs {
				matches, err := doublestar.Glob(glob)
				if err != nil {
					return fmt.Errorf("bad --query-glob pattern: %s", glob) // ignore err, it's not helpful
				}
				queries = append(queries, matches...)
			}
			for i, query := range queries {
				abs, err := filepath.Abs(query)
				if err != nil {
					return fmt.Errorf("absolute path for %s: %w", query, err)
				}
				queries[i] = abs
			}

			// Get absolute paths for all files matching --schema-glob. Order files
			// lexicographically within each glob but not across all globs.
			schemas := make([]string, 0, len(*schemaGlobs)*4)
			for _, glob := range *schemaGlobs {
				matches, err := doublestar.Glob(glob)
				if err != nil {
					return fmt.Errorf("bad --schema-glob pattern: %s", glob) // ignore err, it's not helpful
				}
				sort.Strings(matches)
				schemas = append(schemas, matches...)
			}
			for i, schema := range schemas {
				abs, err := filepath.Abs(schema)
				if err != nil {
					return fmt.Errorf("absolute path for %s: %w", schema, err)
				}
				schemas[i] = abs
			}

			// Deduce output directory.
			outDir := *outputDir
			if outDir == "" {
				for _, file := range queries {
					dir := filepath.Dir(file)
					if outDir != "" && dir != outDir {
						return fmt.Errorf("cannot deduce output dir because query files use different dirs; " +
							"specify explicitly with --output-dir")
					}
					outDir = dir
				}
			}
			// Codegen.
			err := pggen.Generate(pggen.GenerateOptions{
				Language:          pggen.LangGo,
				ConnString:        *postgresConn,
				DockerInitScripts: schemas,
				QueryFiles:        queries,
				OutputDir:         outDir,
			})
			fmt.Printf("gen go: out_dir=%s files=%s\n", outDir, strings.Join(queries, ","))
			return err
		},
	}
	cmd := &ffcli.Command{
		Name:        "gen",
		ShortUsage:  "pggen gen (go|<lang>) [options...]",
		ShortHelp:   "generates code in specific language for Postgres query files",
		FlagSet:     nil,
		Subcommands: []*ffcli.Command{goSubCmd},
	}
	cmd.Exec = func(ctx context.Context, args []string) error {
		fmt.Println(ffcli.DefaultUsageFunc(cmd))
		os.Exit(1)
		return nil
	}
	return cmd
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
