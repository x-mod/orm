package object

import (
	"strings"

	"github.com/x-mod/errors"
)

type FieldD struct {
	Name       string                 `yaml:"name,omitempty"`
	Type       string                 `yaml:"type,omitempty"`
	Comment    string                 `yaml:"comment,omitempty"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
	Tags       map[string]interface{} `yaml:"tags,omitempty"`
}

type UniqueD struct {
	Name       string   `yaml:"name,omitempty"`
	FieldNames []string `yaml:"fields,omitempty"`
}

type IndexD struct {
	Name       string   `yaml:"name,omitempty"`
	FieldNames []string `yaml:"fields,omitempty"`
}

type Object struct {
	Version    string                 `yaml:"version,omitempty"`
	Object     string                 `yaml:"object,omitempty"`
	Name       string                 `yaml:"name,omitempty"`
	Comment    string                 `yaml:"comment,omitempty"`
	Fields     []*FieldD              `yaml:"fields,omitempty"`
	Uniques    []*UniqueD             `yaml:"uniques,omitempty"`
	Indexes    []*IndexD              `yaml:"indexes,omitempty"`
	Primary    []string               `yaml:"primary,omitempty"`
	Tags       []string               `yaml:"tags,omitempty"`
	Attributes map[string]interface{} `yaml:"attributes,omitempty"`
}

func (obj *Object) Adjust() {
	obj.Comment = strings.ReplaceAll(obj.Comment, ";", "")
	for _, field := range obj.Fields {
		field.Comment = strings.ReplaceAll(field.Comment, ";", "")
	}
}
func (obj *Object) IsTable() bool {
	if strings.ToLower(obj.Object) == "table" {
		return true
	}
	return false
}

func (obj *Object) IsQuery() bool {
	if strings.ToLower(obj.Object) == "query" {
		return true
	}
	return false
}

func (obj *Object) Table() (*Table, error) {
	if strings.ToLower(obj.Object) == "table" {
		if len(obj.Primary) == 0 {
			return nil, errors.Errorf("table (%s) primary key required", obj.Name)
		}
		tb := &Table{}
		tb.name = obj.Name
		tb.comment = obj.Comment
		tb.attributes = obj.Attributes
		tb.tags = obj.Tags
		fs := make([]*Field, 0, len(obj.Fields))
		for _, fd := range obj.Fields {
			if !IsFieldType(fd.Type) {
				return nil, errors.Errorf("table (%s) field type (%s) unsupport", obj.Name, fd.Type)
			}
			field := newField(fd)
			fs = append(fs, field)
		}
		tb.fields = fs
		tb.uniques = make([]*Unique, 0, len(obj.Uniques))
		for _, uk := range obj.Uniques {
			unique := &Unique{tb: tb, name: strings.ToLower(uk.Name)}
			unique.fields = make([]*Field, 0, len(uk.FieldNames))
			for _, n := range uk.FieldNames {
				field := tb.FieldByName(n)
				if field == nil {
					return nil, errors.Errorf("table (%s) unique field (%s) not exist", tb.name, n)
				}
				unique.fields = append(unique.fields, field)
			}
			tb.uniques = append(tb.uniques, unique)
		}
		tb.indexes = make([]*Index, 0, len(obj.Indexes))
		for _, idx := range obj.Indexes {
			index := &Index{tb: tb, name: strings.ToLower(idx.Name)}
			index.fields = make([]*Field, 0, len(idx.FieldNames))
			for _, n := range idx.FieldNames {
				field := tb.FieldByName(n)
				if field == nil {
					return nil, errors.Errorf("table (%s) index field (%s) not exist", tb.name, n)
				}
				index.fields = append(index.fields, field)
			}
			tb.indexes = append(tb.indexes, index)
		}
		primary := &Primary{tb: tb}
		primary.fields = make([]*Field, 0, len(obj.Primary))
		for _, n := range obj.Primary {
			field := tb.FieldByName(n)
			if field == nil {
				return nil, errors.Errorf("table (%s) index field (%s) not exist", tb.name, n)
			}
			field.primary = true
			primary.fields = append(primary.fields, field)
		}
		tb.primary = primary
		return tb, nil
	}
	return nil, errors.Errorf("object (%s) is not table", obj.Name)
}

func (obj *Object) Query() (*Query, error) {
	if strings.ToLower(obj.Object) == "query" {
		qry := &Query{}
		qry.name = obj.Name
		qry.tags = obj.Tags
		fs := make([]*Field, 0, len(obj.Fields))
		fm := make(map[string]*Field)
		for _, fd := range obj.Fields {
			if !IsFieldType(fd.Type) {
				return nil, errors.Errorf("query (%s) field type (%s) unsupport", obj.Name, fd.Type)
			}
			field := newField(fd)
			fs = append(fs, field)
			fm[field.name] = field
		}
		qry.fields = fs
		qry.mapFields = fm
		return qry, nil
	}
	return nil, errors.Errorf("object (%s) is not query", obj.Name)
}
