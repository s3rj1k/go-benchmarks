package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"testing"
	"text/template"

	sq "github.com/Masterminds/squirrel"
	"github.com/cristalhq/builq"
	"github.com/flosch/pongo2/v6"
	"github.com/keegancsmith/sqlf"
	_ "github.com/mattn/go-sqlite3"
)

// go test -bench=. -benchmem

func BenchmarkSQLiteInsertSelectUpdate(b *testing.B) { // dynamic, unsafe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
		id INTEGER PRIMARY KEY,
		name TEXT,
		value1 REAL,
		value2 REAL,
		value3 REAL,
		value4 REAL,
		value5 REAL,
		value6 REAL,
		value7 REAL,
		value8 REAL,
		value9 REAL
	)`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.Exec("INSERT INTO benchmark(name, value1, value2, value3, value4, value5, value6, value7, value8, value9) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			fmt.Sprintf("Name%d", i), float64(i), float64(i+1), float64(i+2), float64(i+3), float64(i+4), float64(i+5), float64(i+6), float64(i+7), float64(i+8))
		if err != nil {
			b.Fatalf("could not execute statement: %v", err)
		}

		var id int
		err = db.QueryRow("SELECT id FROM benchmark WHERE name = ?", fmt.Sprintf("Name%d", i)).Scan(&id)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatalf("data mismatch: expected %d, got %d.", i+1, id)
		}

		_, err = db.Exec("UPDATE benchmark SET value1 = ? WHERE id = ?", float64(i+10), id)
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		var val float64
		err = db.QueryRow("SELECT value1 FROM benchmark WHERE id = ?", id).Scan(&val)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingSquirrel(b *testing.B) { // typed, safe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
		id INTEGER PRIMARY KEY,
		name TEXT,
		value1 REAL,
		value2 REAL,
		value3 REAL,
		value4 REAL,
		value5 REAL,
		value6 REAL,
		value7 REAL,
		value8 REAL,
		value9 REAL
	)`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query, args, err := psql.Insert("benchmark").Columns("name", "value1", "value2", "value3", "value4", "value5", "value6", "value7", "value8", "value9").
			Values(fmt.Sprintf("Name%d", i), float64(i), float64(i+1), float64(i+2), float64(i+3), float64(i+4), float64(i+5), float64(i+6), float64(i+7), float64(i+8)).ToSql()
		if err != nil {
			b.Fatalf("could not build insert SQL: %v", err)
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		query, args, err = psql.Select("id").From("benchmark").Where(sq.Eq{"name": fmt.Sprintf("Name%d", i)}).ToSql()
		if err != nil {
			b.Fatalf("could not build select SQL: %v", err)
		}

		var id int
		err = db.QueryRow(query, args...).Scan(&id)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatalf("data mismatch: expected %d, got %d.", i+1, id)
		}

		query, args, err = psql.Update("benchmark").Set("value1", float64(i+10)).Where(sq.Eq{"id": id}).ToSql()
		if err != nil {
			b.Fatalf("could not build update SQL: %v", err)
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		query, args, err = psql.Select("value1").From("benchmark").Where(sq.Eq{"id": id}).ToSql()
		if err != nil {
			b.Fatalf("could not build select SQL: %v", err)
		}

		var val float64
		err = db.QueryRow(query, args...).Scan(&val)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingSqlf(b *testing.B) { // semi-dynamic, safe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
		id INTEGER PRIMARY KEY,
		name TEXT,
		value1 REAL,
		value2 REAL,
		value3 REAL,
		value4 REAL,
		value5 REAL,
		value6 REAL,
		value7 REAL,
		value8 REAL,
		value9 REAL
	)`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		query := sqlf.Sprintf("INSERT INTO benchmark(name, value1, value2, value3, value4, value5, value6, value7, value8, value9) VALUES(%s, %f, %f, %f, %f, %f, %f, %f, %f, %f)",
			fmt.Sprintf("Name%d", i), float64(i), float64(i+1), float64(i+2), float64(i+3), float64(i+4), float64(i+5), float64(i+6), float64(i+7), float64(i+8))

		_, err := db.Exec(query.Query(sqlf.SQLServerBindVar), query.Args()...)
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		query = sqlf.Sprintf("SELECT id FROM benchmark WHERE name = %s", fmt.Sprintf("Name%d", i))

		var id int
		err = db.QueryRow(query.Query(sqlf.SQLServerBindVar), query.Args()...).Scan(&id)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatalf("data mismatch: expected %d, got %d.", i+1, id)
		}

		query = sqlf.Sprintf("UPDATE benchmark SET value1 = %f WHERE id = %d", float64(i+10), id)

		_, err = db.Exec(query.Query(sqlf.SQLServerBindVar), query.Args()...)
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		query = sqlf.Sprintf("SELECT value1 FROM benchmark WHERE id = %s", id)

		var val float64
		err = db.QueryRow(query.Query(sqlf.SQLServerBindVar), query.Args()...).Scan(&val)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingTemplateWithMap(b *testing.B) { // unsafe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
        id INTEGER PRIMARY KEY,
        name TEXT,
        value1 REAL,
        value2 REAL,
        value3 REAL,
        value4 REAL,
        value5 REAL,
        value6 REAL,
        value7 REAL,
        value8 REAL,
        value9 REAL
    )`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	insertTemplate := `INSERT INTO benchmark(name, value1, value2, value3, value4, value5, value6, value7, value8, value9) VALUES('{{.Name}}', {{.Value1}}, {{.Value2}}, {{.Value3}}, {{.Value4}}, {{.Value5}}, {{.Value6}}, {{.Value7}}, {{.Value8}}, {{.Value9}})`
	selectIdTemplate := `SELECT id FROM benchmark WHERE name = '{{.Name}}'`
	updateTemplate := `UPDATE benchmark SET value1 = {{.Value1}} WHERE id = {{.Id}}`
	selectValueTemplate := `SELECT value1 FROM benchmark WHERE id = {{.Id}}`

	tmplInsert := template.Must(template.New("insert").Parse(insertTemplate))
	tmplSelectId := template.Must(template.New("selectId").Parse(selectIdTemplate))
	tmplUpdate := template.Must(template.New("update").Parse(updateTemplate))
	tmplSelectValue := template.Must(template.New("selectValue").Parse(selectValueTemplate))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer

		err = tmplInsert.Execute(&buf, map[string]any{
			"Name":   fmt.Sprintf("Name%d", i),
			"Value1": float64(i),
			"Value2": float64(i + 1),
			"Value3": float64(i + 2),
			"Value4": float64(i + 3),
			"Value5": float64(i + 4),
			"Value6": float64(i + 5),
			"Value7": float64(i + 6),
			"Value8": float64(i + 7),
			"Value9": float64(i + 8),
		})
		if err != nil {
			b.Fatalf("could not execute insert template: %v", err)
		}

		_, err = db.Exec(buf.String())
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		err = tmplSelectId.Execute(&buf, map[string]any{
			"Name": fmt.Sprintf("Name%d", i),
		})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var id int
		err = db.QueryRow(buf.String()).Scan(&id)
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatalf("data mismatch: expected %d, got %d.", i+1, id)
		}

		err = tmplUpdate.Execute(&buf, map[string]any{
			"Value1": float64(i + 10),
			"Id":     id,
		})
		if err != nil {
			b.Fatalf("could not execute update template: %v", err)
		}

		_, err = db.Exec(buf.String())
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		err = tmplSelectValue.Execute(&buf, map[string]any{
			"Id": id,
		})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var val float64
		err = db.QueryRow(buf.String()).Scan(&val)
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingTemplateWithStruct(b *testing.B) { // unsafe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
        id INTEGER PRIMARY KEY,
        name TEXT,
        value1 REAL,
        value2 REAL,
        value3 REAL,
        value4 REAL,
        value5 REAL,
        value6 REAL,
        value7 REAL,
        value8 REAL,
        value9 REAL
    )`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	insertTemplate := `INSERT INTO benchmark(name, value1, value2, value3, value4, value5, value6, value7, value8, value9) VALUES('{{.Name}}', {{.Value1}}, {{.Value2}}, {{.Value3}}, {{.Value4}}, {{.Value5}}, {{.Value6}}, {{.Value7}}, {{.Value8}}, {{.Value9}})`
	selectIdTemplate := `SELECT id FROM benchmark WHERE name = '{{.Name}}'`
	updateTemplate := `UPDATE benchmark SET value1 = {{.Value1}} WHERE id = {{.Id}}`
	selectValueTemplate := `SELECT value1 FROM benchmark WHERE id = {{.Id}}`

	tmplInsert := template.Must(template.New("insert").Parse(insertTemplate))
	tmplSelectId := template.Must(template.New("selectId").Parse(selectIdTemplate))
	tmplUpdate := template.Must(template.New("update").Parse(updateTemplate))
	tmplSelectValue := template.Must(template.New("selectValue").Parse(selectValueTemplate))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer

		err = tmplInsert.Execute(&buf, struct {
			Name   string
			Value1 float64
			Value2 float64
			Value3 float64
			Value4 float64
			Value5 float64
			Value6 float64
			Value7 float64
			Value8 float64
			Value9 float64
		}{
			Name:   fmt.Sprintf("Name%d", i),
			Value1: float64(i),
			Value2: float64(i + 1),
			Value3: float64(i + 2),
			Value4: float64(i + 3),
			Value5: float64(i + 4),
			Value6: float64(i + 5),
			Value7: float64(i + 6),
			Value8: float64(i + 7),
			Value9: float64(i + 8),
		})
		if err != nil {
			b.Fatalf("could not execute insert template: %v", err)
		}

		_, err = db.Exec(buf.String())
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		err = tmplSelectId.Execute(&buf, struct {
			Name string
		}{
			Name: fmt.Sprintf("Name%d", i),
		})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var id int
		err = db.QueryRow(buf.String()).Scan(&id)
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatal("invalid data coming from DB")
		}

		err = tmplUpdate.Execute(&buf, struct {
			Value1 float64
			Id     int
		}{
			Value1: float64(i + 10),
			Id:     id,
		})
		if err != nil {
			b.Fatalf("could not execute update template: %v", err)
		}

		_, err = db.Exec(buf.String())
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		err = tmplSelectValue.Execute(&buf, struct {
			Id int
		}{
			Id: id,
		})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var val float64
		err = db.QueryRow(buf.String()).Scan(&val)
		buf.Reset()
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingPongo2(b *testing.B) { // unsafe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
        id INTEGER PRIMARY KEY,
        name TEXT,
        value1 REAL,
        value2 REAL,
        value3 REAL,
        value4 REAL,
        value5 REAL,
        value6 REAL,
        value7 REAL,
        value8 REAL,
        value9 REAL
    )`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	insertTplString := "INSERT INTO benchmark(name, value1, value2, value3, value4, value5, value6, value7, value8, value9) VALUES('{{ name }}', {{ value1 }}, {{ value2 }}, {{ value3 }}, {{ value4 }}, {{ value5 }}, {{ value6 }}, {{ value7 }}, {{ value8 }}, {{ value9 }})"
	selectIdTplString := "SELECT id FROM benchmark WHERE name = '{{ name }}'"
	updateTplString := "UPDATE benchmark SET value1 = {{ value1 }} WHERE id = {{ id }}"
	selectValueTplString := "SELECT value1 FROM benchmark WHERE id = {{ id }}"

	insertTpl, err := pongo2.FromString(insertTplString)
	if err != nil {
		b.Fatalf("could not parse insert template: %v", err)
	}

	selectIdTpl, err := pongo2.FromString(selectIdTplString)
	if err != nil {
		b.Fatalf("could not parse select template: %v", err)
	}

	updateTpl, err := pongo2.FromString(updateTplString)
	if err != nil {
		b.Fatalf("could not parse update template: %v", err)
	}

	selectValueTpl, err := pongo2.FromString(selectValueTplString)
	if err != nil {
		b.Fatalf("could not parse select template: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sql, err := insertTpl.Execute(pongo2.Context{
			"name":   fmt.Sprintf("Name%d", i),
			"value1": float64(i),
			"value2": float64(i + 1),
			"value3": float64(i + 2),
			"value4": float64(i + 3),
			"value5": float64(i + 4),
			"value6": float64(i + 5),
			"value7": float64(i + 6),
			"value8": float64(i + 7),
			"value9": float64(i + 8),
		})
		if err != nil {
			b.Fatalf("could not execute insert template: %v", err)
		}

		_, err = db.Exec(sql)
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		sql, err = selectIdTpl.Execute(pongo2.Context{"name": fmt.Sprintf("Name%d", i)})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var id int
		err = db.QueryRow(sql).Scan(&id)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatal("invalid data coming from DB")
		}

		sql, err = updateTpl.Execute(pongo2.Context{"value1": float64(i + 10), "id": id})
		if err != nil {
			b.Fatalf("could not execute update template: %v", err)
		}

		_, err = db.Exec(sql)
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		sql, err = selectValueTpl.Execute(pongo2.Context{"id": id})
		if err != nil {
			b.Fatalf("could not execute select template: %v", err)
		}

		var val float64
		err = db.QueryRow(sql).Scan(&val)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}

func BenchmarkSQLiteInsertSelectUpdateUsingBuilq(b *testing.B) { // dynamic, safe SQL
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		b.Fatalf("could not open sqlite3 database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE benchmark (
        id INTEGER PRIMARY KEY,
        name TEXT,
        value1 REAL,
        value2 REAL,
        value3 REAL,
        value4 REAL,
        value5 REAL,
        value6 REAL,
        value7 REAL,
        value8 REAL,
        value9 REAL
    )`)
	if err != nil {
		b.Fatalf("could not create table: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bb := builq.Builder{}

		bb.Addf("INSERT INTO benchmark (%s)",
			builq.Columns{"name", "value1", "value2", "value3", "value4", "value5", "value6", "value7", "value8", "value9"},
		)
		bb.Addf("VALUES (%+$)",
			[]any{fmt.Sprintf("Name%d", i), float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4), float64(i + 5), float64(i + 6), float64(i + 7), float64(i + 8)},
		)

		query, args, err := bb.Build()
		if err != nil {
			b.Fatalf("could not build insert query: %v", err)
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			b.Fatalf("could not execute insert statement: %v", err)
		}

		bf := builq.New()
		bf("SELECT %s FROM %s", "id", "benchmark")
		bf("WHERE %s = %$", "name", fmt.Sprintf("Name%d", i))

		query, args, err = bf.Build()
		if err != nil {
			b.Fatalf("could not build select query: %v", err)
		}

		var id int
		err = db.QueryRow(query, args...).Scan(&id)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if id != (i + 1) {
			b.Fatalf("data mismatch: expected %d, got %d.", i+1, id)
		}

		bb = builq.Builder{}
		bb.Addf("UPDATE benchmark SET value1 = %$ WHERE id = %$", float64(i+10), id)

		query, args, err = bb.Build()
		if err != nil {
			b.Fatalf("could not build update query: %v", err)
		}

		_, err = db.Exec(query, args...)
		if err != nil {
			b.Fatalf("could not execute update statement: %v", err)
		}

		bf = builq.New()
		bf("SELECT %s FROM %s", "value1", "benchmark")
		bf("WHERE %s = %$", "id", id)

		query, args, err = bf.Build()
		if err != nil {
			b.Fatalf("could not build select query: %v", err)
		}

		var val float64
		err = db.QueryRow(query, args...).Scan(&val)
		if err != nil {
			b.Fatalf("could not execute select statement: %v", err)
		}

		if val != float64(i+10) {
			b.Fatalf("data mismatch: expected %f, got %f.", float64(i+10), val)
		}
	}
}
