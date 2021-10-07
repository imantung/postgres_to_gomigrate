package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	DBUser   = "user"
	DBPass   = "pass"
	DBHost   = "localhost"
	DBPort   = "5434"
	DBName   = "user"
	DBSchema = "public"

	TargetFolder = "migrations"

	PgDumpArgs = []string{
		"--no-comments",
		"--no-publications",
		"--no-security-labels",
		"--no-subscriptions",
		"--no-synchronized-snapshots",
		"--no-tablespaces",
		"--no-unlogged-table-data",
		"--no-owner",
		"--no-privileges",
		"--no-blobs",
		"--schema-only",
		"--clean",
	}

	SkipTables = []string{
		"schema_migrations",
	}
)

func main() {
	os.Setenv("PGPASSWORD", DBPass)
	os.Mkdir(TargetFolder, 0777)

	tables, err := tableList(DBSchema)
	if err != nil {
		log.Fatal(err)
	}

	for i, table := range tables {
		log.Printf("Generate migration for '%s'", table)
		err = generateMigrations(strconv.Itoa(i+1), table)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateMigrations(version, table string) error {
	b, err := pgDump(table).CombinedOutput()
	if err != nil {
		return errors.New(string(b))
	}

	var downLines []string
	var upLines []string
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if isSkipLine(line) {
			// do nothing
		} else if isDownScript(line) {
			downLines = append(downLines, line)
		} else {
			upLines = append(upLines, line)
		}
	}

	filename := fmt.Sprintf("%s/%s_%s", TargetFolder, version, table)
	dumpFilename := filename + ".dump.sql"
	upFilename := filename + ".up.sql"
	downFilename := filename + ".down.sql"

	if err := ioutil.WriteFile(dumpFilename, b, 0777); err != nil {
		return err
	}

	downScript := []byte(strings.Join(downLines, "\n"))
	if err := ioutil.WriteFile(downFilename, downScript, 0777); err != nil {
		return err
	}

	upScript := []byte(strings.Join(upLines, "\n"))
	return ioutil.WriteFile(upFilename, upScript, 0777)
}

func isSkipLine(line string) bool {
	return line == "" ||
		strings.HasPrefix(line, "--") ||
		strings.HasPrefix(line, "SET ") ||
		strings.HasPrefix(line, "SELECT ")
}

func isDownScript(line string) bool {
	return strings.Contains(line, "DROP")
}

func pgDump(table string) *exec.Cmd {
	var args []string
	args = append(PgDumpArgs,
		"--username", DBUser,
		"--port", DBPort,
		"--host", DBHost,
		"--table", table,
		DBName,
	)
	return exec.Command("pg_dump", args...)
}

func isSkipTable(table string) bool {
	for _, skipTable := range SkipTables {
		if skipTable == table {
			return true
		}
	}
	return false
}

func tableList(schema string) ([]string, error) {
	queryFormat := `SELECT table_name FROM information_schema.tables WHERE table_schema='%s' AND table_type='BASE TABLE';`
	query := fmt.Sprintf(queryFormat, schema)
	b, err := pgSQL(query).CombinedOutput()
	if err != nil {
		return nil, errors.New(string(b))
	}

	tables := strings.Split(string(b), "\n")
	if len(tables) < 5 {
		return nil, nil
	}

	tables = tables[2 : len(tables)-3]
	var result []string
	for _, table := range tables {
		table = strings.TrimSpace(table)
		if !isSkipTable(table) {
			result = append(result, table)
		}
	}
	return result, nil
}

func pgSQL(query string) *exec.Cmd {
	return exec.Command("psql",
		"-h", DBHost,
		"-p", DBPort,
		"-U", DBUser,
		"-d", DBName,
		"-c", query,
	)
}
