package object

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/stoewer/go-strcase"
)

type Field struct {
	name       string
	typeName   string
	comment    string
	attributes map[string]interface{}
	tags       map[string]interface{}
	primary    bool
}

func newField(fd *FieldD) *Field {
	return &Field{
		name:       fd.Name,
		typeName:   fd.Type,
		comment:    fd.Comment,
		attributes: fd.Attributes,
		tags:       fd.Tags,
	}
}

func (f *Field) Name() string {
	return strcase.UpperCamelCase(f.name)
}

func (f *Field) Value() string {
	if f.IsTime() {
		return fmt.Sprintf("%s.Unix()", f.Name())
	}
	return f.Name()
}

func (f *Field) Type(language string) string {
	//TODO type convert
	switch strings.ToLower(language) {
	case "golang":
		t, ok := fieldGoTypes[f.typeName]
		if !ok {
			return "string"
		}
		if t == "time.Time" && f.IsNullable() {
			return "*time.Time"
		}
		return t
	case "db":
		if c := f.Attribute("dbType"); c != nil {
			return cast.ToString(c)
		}
		if f.IsText() {
			return "TEXT"
		}
		return DBType(f.typeName)
	case "mysql":
		if c := f.Attribute("mysqlType"); c != nil {
			return cast.ToString(c)
		}
		if c := f.Attribute("dbType"); c != nil {
			return cast.ToString(c)
		}
		if f.IsText() {
			return "TEXT"
		}
		return MySQLType(f.typeName)
	}
	return f.typeName
}

func (f *Field) SQLNullValueField() string {
	gotype := f.Type("golang")
	if f.IsTime() {
		return "Int64"
	}
	if f.IsNullable() {
		if strings.Contains(gotype, "bool") {
			return "Bool"
		}
		if strings.Contains(gotype, "string") {
			return "String"
		}
		if strings.Contains(gotype, "int64") {
			return "Int64"
		}
		if strings.Contains(gotype, "int32") {
			return "Int32"
		}
		if strings.Contains(gotype, "int") {
			return "Int32"
		}
		if strings.Contains(gotype, "float") {
			return "Float64"
		}
		if strings.Contains(gotype, "time.Time") {
			return "Int64"
		}
	}
	return gotype
}
func (f *Field) SQLNullType() string {
	gotype := f.Type("golang")
	if f.IsTime() {
		return "sql.NullInt64"
	}
	if f.IsNullable() {
		if strings.Contains(gotype, "bool") {
			return "sql.NullBool"
		}
		if strings.Contains(gotype, "string") {
			return "sql.NullString"
		}
		if strings.Contains(gotype, "int64") {
			return "sql.NullInt64"
		}
		if strings.Contains(gotype, "int32") {
			return "sql.NullInt32"
		}
		if strings.Contains(gotype, "int") {
			return "sql.NullInt32"
		}
		if strings.Contains(gotype, "float") {
			return "sql.NullFloat64"
		}
		if strings.Contains(gotype, "time.Time") {
			return "sql.NullInt64"
		}
	}
	return gotype
}

func (f *Field) IsTime() bool {
	gotype := f.Type("golang")
	if strings.Contains(gotype, "time.Time") {
		return true
	}
	return false
}

func (f *Field) Comment() string {
	return f.comment
}

func (f *Field) DBType() string {
	return f.Type("db")
}
func (f *Field) DBColumn() string {
	if c := f.Attribute("dbColumn"); c != nil {
		return "`" + strcase.SnakeCase(cast.ToString(c)) + "`"
	}

	return "`" + strcase.SnakeCase(f.name) + "`"
}
func (f *Field) DBDefault() string {
	if f.IsAutoIncrement() {
		return "AUTO_INCREMENT"
	}
	if f.IsText() {
		return ""
	}
	if c := f.Attribute("dbDefault"); c != nil {
		return cast.ToString(c)
	}
	gotype := f.Type("golang")
	if strings.Contains(gotype, "bool") {
		return "DEFAULT \"0\""
	}
	if strings.Contains(gotype, "string") {
		return "DEFAULT \"\""
	}
	if strings.Contains(gotype, "int") {
		return "DEFAULT \"0\""
	}
	if strings.Contains(gotype, "float") {
		return "DEFAULT \"0\""
	}
	if strings.Contains(gotype, "time.Time") {
		return "DEFAULT \"0\""
	}
	return ""
}

func (f *Field) DBNull() string {
	if f.IsNullable() {
		return "NULL"
	}
	return "NOT NULL"
}

func (f *Field) IsNullable() bool {
	if c := f.Attribute("nullable"); c != nil {
		return cast.ToBool(c)
	}
	return false
}

func (f *Field) IsPrimary() bool {
	return f.primary
}

func (f *Field) IsAutoIncrement() bool {
	if c := f.Attribute("autoIncrement"); c != nil {
		return cast.ToBool(c)
	}
	return false
}

func (f *Field) IsText() bool {
	if c := f.Attribute("text"); c != nil {
		return cast.ToBool(c)
	}
	return false
}

func (f *Field) IsPassword() bool {
	if c := f.Attribute("password"); c != nil {
		return cast.ToBool(c)
	}
	return false
}

func (f *Field) IsEncode() bool {
	if c := f.Attribute("encode"); c != nil {
		return cast.ToBool(c)
	}
	return false
}

func (f *Field) TagValue(key string) string {
	if v, ok := f.tags[key]; ok {
		return cast.ToString(v)
	}
	return strcase.LowerCamelCase(f.name)
}

func (f *Field) Tags(keys []string) string {
	tags := []string{}
	for _, k := range keys {
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", k, f.TagValue(k)))
	}
	return strings.Join(tags, " ")
}

func (f *Field) Attribute(name string) interface{} {
	if v, ok := f.attributes[name]; ok {
		return v
	}
	return nil
}
