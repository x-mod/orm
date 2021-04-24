package object

import (
	"strings"

	"github.com/spf13/cast"

	"github.com/stoewer/go-strcase"
)

// "github.com/spf13/cast"

type Unique struct {
	tb      *Table
	name    string
	comment string
	fields  []*Field
}

func (u *Unique) Name() string {
	if len(u.name) == 0 {
		names := make([]string, 0, len(u.fields))
		for _, f := range u.fields {
			names = append(names, f.Name())
		}
		return u.tb.Name() + strings.Join(names, "")
	}
	return u.tb.Name() + strcase.UpperCamelCase(u.name)
}

func (u *Unique) Table() *Table {
	return u.tb
}

func (u *Unique) Comment() string {
	return u.comment
}

func (u *Unique) Fields() []*Field {
	return u.fields
}

func (u *Unique) JoinFields(sep string) string {
	names := []string{}
	for _, f := range u.fields {
		names = append(names, f.DBColumn())
	}
	return strings.Join(names, sep)
}

type Index struct {
	tb      *Table
	name    string
	comment string
	fields  []*Field
}

func (u *Index) Name() string {
	if len(u.name) == 0 {
		names := make([]string, 0, len(u.fields))
		for _, f := range u.fields {
			names = append(names, f.Name())
		}
		return u.tb.Name() + strings.Join(names, "")
	}
	return u.tb.Name() + strcase.UpperCamelCase(u.name)
}

func (u *Index) Table() *Table {
	return u.tb
}

func (u *Index) Comment() string {
	return u.comment
}

func (u *Index) Fields() []*Field {
	return u.fields
}

func (u *Index) JoinFields(sep string) string {
	names := []string{}
	for _, f := range u.fields {
		names = append(names, f.DBColumn())
	}
	return strings.Join(names, sep)
}

type Primary struct {
	tb     *Table
	fields []*Field
}

func (u *Primary) Table() *Table {
	return u.tb
}

func (u *Primary) Fields() []*Field {
	return u.fields
}

func (u *Primary) JoinFields(sep string) string {
	names := []string{}
	for _, f := range u.fields {
		names = append(names, f.DBColumn())
	}
	return strings.Join(names, sep)
}

type Table struct {
	name       string
	comment    string
	fields     []*Field
	uniques    []*Unique
	indexes    []*Index
	primary    *Primary
	attributes map[string]interface{}
	tags       []string
}

func (tb *Table) Name() string {
	return strcase.UpperCamelCase(tb.name)
}

func (tb *Table) DBName() string {
	if v, ok := tb.attributes["table"]; ok {
		return cast.ToString(v)
	}
	return strcase.SnakeCase(tb.name)
}

func (tb *Table) DBEngine(defaultEngine string) string {
	if v, ok := tb.attributes["engine"]; ok {
		return cast.ToString(v)
	}
	return defaultEngine
}

func (tb *Table) DBCharset(defaultCharset string) string {
	if v, ok := tb.attributes["charset"]; ok {
		return cast.ToString(v)
	}
	return defaultCharset
}

func (tb *Table) AutoIncrement(defautBegin string) string {
	if v, ok := tb.attributes["autoIncrement"]; ok {
		return cast.ToString(v)
	}
	return defautBegin
}

func (tb *Table) Comment() string {
	return tb.comment
}

func (tb *Table) Attribute(name string) interface{} {
	if v, ok := tb.attributes[name]; ok {
		return v
	}
	return nil
}

func (tb *Table) Tags() []string {
	return tb.tags
}

func (tb *Table) Fields() []*Field {
	return tb.fields
}

func (tb *Table) FieldByName(name string) *Field {
	for _, f := range tb.fields {
		if strcase.UpperCamelCase(f.name) == strcase.UpperCamelCase(name) {
			return f
		}
	}
	return nil
}

func (tb *Table) Uniques() []*Unique {
	return tb.uniques
}

func (tb *Table) UniqueByName(name string) *Unique {
	for _, f := range tb.uniques {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (tb *Table) Indexes() []*Index {
	return tb.indexes
}

func (tb *Table) IndexByName(name string) *Index {
	for _, f := range tb.indexes {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (tb *Table) PrimaryKey() *Primary {
	return tb.primary
}

type Query struct {
	name      string
	fields    []*Field
	mapFields map[string]*Field
	tags      []string
}

func (q *Query) Name() string {
	return strcase.UpperCamelCase(q.name)
}

func (q *Query) Fields() []*Field {
	return q.fields
}

func (q *Query) FieldByName(name string) *Field {
	for _, f := range q.fields {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (q *Query) Tags() []string {
	return q.tags
}
