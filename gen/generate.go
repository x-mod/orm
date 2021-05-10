package gen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/x-mod/errors"
	"github.com/x-mod/orm/object"
	"github.com/x-mod/orm/parser"
)

func generateScripts(in string, out string, suffix string, objects ...string) error {
	t, err := getTemplate(strings.Join(objects, "/"))
	if err != nil {
		return err
	}

	inputs := []string{}
	for _, sfx := range inputSuffixes {
		files, err := getInputFilesBySuffix(in, sfx)
		if err != nil {
			return errors.Annotate(err, "get input files")
		}
		inputs = append(inputs, files...)
	}

	tables := []*object.Table{}
	for _, input := range inputs {
		objs, err := parser.Parse(input)
		if err != nil {
			return err
		}
		for _, obj := range objs {
			if obj.IsTable() {
				tb, err := obj.Table()
				if err != nil {
					return err
				}
				tables = append(tables, tb)
			}
		}
	}

	for _, tb := range tables {
		outfile := filepath.Join(out, fmt.Sprintf("%s.sql", tb.DBName()))
		wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Annotatef(err, "write file %s", outfile)
		}
		if err := t.ExecuteTemplate(wr, "table", map[string]interface{}{
			"PackageName": packageName,
			"Table":       tb,
		}); err != nil {
			return errors.Annotatef(err, "template execute %s", tb.Name())
		}
		if err := wr.Close(); err != nil {
			return errors.Annotatef(err, "close file %s", outfile)
		}
	}
	return copyFilesExcludeSuffix(strings.Join(objects, "/"), suffix, out, true)
}

func generateCode(in string, out string, suffix string, objects ...string) error {
	t, err := getTemplate(strings.Join(objects, "/"))
	if err != nil {
		return err
	}

	inputs := []string{}
	for _, sfx := range inputSuffixes {
		files, err := getInputFilesBySuffix(in, sfx)
		if err != nil {
			return errors.Annotate(err, "get input files")
		}
		inputs = append(inputs, files...)
	}

	tables := []*object.Table{}
	queries := []*object.Query{}
	for _, input := range inputs {
		objs, err := parser.Parse(input)
		if err != nil {
			return err
		}
		for _, obj := range objs {
			if obj.IsTable() {
				tb, err := obj.Table()
				if err != nil {
					return err
				}
				tables = append(tables, tb)
			}
			if obj.IsQuery() {
				q, err := obj.Query()
				if err != nil {
					return err
				}
				queries = append(queries, q)
			}
		}
	}

	//db
	if true {
		outfile := filepath.Join(out, fmt.Sprintf("gen.db.go"))
		wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Annotatef(err, "write file %s", outfile)
		}
		if err := t.ExecuteTemplate(wr, "db", map[string]interface{}{
			"PackageName": packageName,
		}); err != nil {
			return errors.Annotatef(err, "template execute db")
		}
		if err := wr.Close(); err != nil {
			return errors.Annotatef(err, "close file %s", outfile)
		}
		oscmd := exec.Command("gofmt", "-w", outfile)
		if err := oscmd.Run(); err != nil {
			log.Println("gofmt: ", err)
		}
	}
	for _, tb := range tables {
		outfile := filepath.Join(out, fmt.Sprintf("gen.table.%s.go", tb.Name()))
		wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Annotatef(err, "write file %s", outfile)
		}
		if err := t.ExecuteTemplate(wr, "table", map[string]interface{}{
			"PackageName": packageName,
			"Table":       tb,
		}); err != nil {
			return errors.Annotatef(err, "template execute %s", tb.Name())
		}
		if err := wr.Close(); err != nil {
			return errors.Annotatef(err, "close file %s", outfile)
		}
		oscmd := exec.Command("gofmt", "-w", outfile)
		if err := oscmd.Run(); err != nil {
			log.Println("gofmt: ", err)
		}
	}

	for _, q := range queries {
		outfile := filepath.Join(out, fmt.Sprintf("gen.query.%s.go", q.Name()))
		wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return errors.Annotatef(err, "write file %s", outfile)
		}
		if err := t.ExecuteTemplate(wr, "query", map[string]interface{}{
			"PackageName": packageName,
			"Query":       q,
		}); err != nil {
			return errors.Annotatef(err, "template execute %s", q.Name())
		}
		if err := wr.Close(); err != nil {
			return errors.Annotatef(err, "close file %s", outfile)
		}
		oscmd := exec.Command("gofmt", "-w", outfile)
		if err := oscmd.Run(); err != nil {
			log.Println("gofmt: ", err)
		}
	}

	return copyFilesExcludeSuffix(strings.Join(objects, "/"), suffix, out, true)
}
