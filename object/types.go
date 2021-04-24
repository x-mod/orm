package object

var fieldTypes = map[string]bool{
	"bool":      true,
	"int":       true,
	"int8":      true,
	"int16":     true,
	"int32":     true,
	"int64":     true,
	"uint":      true,
	"uint8":     true,
	"uint16":    true,
	"uint32":    true,
	"uint64":    true,
	"float":     true,
	"float32":   true,
	"float64":   true,
	"string":    true,
	"date":      true,
	"time":      true,
	"datetime":  true,
	"timestamp": true,
	"timeint":   true,
}

var fieldGoTypes = map[string]string{
	"bool":      "bool",
	"int":       "int32",
	"int8":      "int8",
	"int16":     "int16",
	"int32":     "int32",
	"int64":     "int64",
	"uint":      "uint32",
	"uint8":     "uint8",
	"uint16":    "uint16",
	"uint32":    "uint32",
	"uint64":    "uint64",
	"float":     "float32",
	"float32":   "float32",
	"float64":   "float64",
	"string":    "string",
	"date":      "time.Time",
	"time":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"timeint":   "time.Time",
}

var fieldMySQLTypes = map[string]string{
	"bool":      "TINYINT(1) UNSIGNED",
	"int":       "INT(11)",
	"int8":      "SMALLINT",
	"int16":     "MEDIUMINT",
	"int32":     "INT(11)",
	"int64":     "BIGINT(20)",
	"uint":      "INT(11) UNSIGNED",
	"uint8":     "SMALLINT UNSIGNED",
	"uint16":    "MEDIUMINT UNSIGNED",
	"uint32":    "INT(11) UNSIGNED",
	"uint64":    "BIGINT(20) UNSIGNED",
	"float":     "FLOAT",
	"float32":   "FLOAT",
	"float64":   "FLOAT",
	"string":    "VARCHAR(128)",
	"date":      "BIGINT(20)",
	"time":      "BIGINT(20)",
	"datetime":  "BIGINT(20)",
	"timestamp": "BIGINT(20)",
	"timeint":   "BIGINT(20)",
}

func IsFieldType(t string) bool {
	if v, ok := fieldTypes[t]; ok {
		return v
	}
	return false
}

func DBType(t string) string {
	return MySQLType(t)
}

func MySQLType(t string) string {
	if v, ok := fieldMySQLTypes[t]; ok {
		return v
	}
	return ""
}
