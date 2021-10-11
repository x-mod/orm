package model

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/x-mod/glog"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

//LOG
type LOGDB struct {
	db DB
}

func LOG(db DB) *LOGDB {
	return &LOGDB{db: db}
}

func (lg *LOGDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	rs, err := lg.db.ExecContext(ctx, query, args...)
	if err != nil {
		glog.Errorf("execute <%s> <%v>: %v", query, args, err)
	} else {
		glog.V(4).Infof("execute <%s> <%v>", query, args)
	}
	return rs, err
}

func (lg *LOGDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rs, err := lg.db.QueryContext(ctx, query, args...)
	if err != nil {
		glog.Errorf("query <%s> <%v>: %v", query, args, err)
	} else {
		glog.V(4).Infof("query <%s> <%v>", query, args)
	}
	return rs, err
}

//In
type FieldIN struct {
	Field   string
	Params  []interface{}
	holders []string
}

func NewFieldIN(field string) *FieldIN {
	return &FieldIN{
		Field:   field,
		Params:  []interface{}{},
		holders: []string{},
	}
}

func (in *FieldIN) Add(v interface{}) *FieldIN {
	in.Params = append(in.Params, v)
	in.holders = append(in.holders, "?")
	return in
}

func (in *FieldIN) SQLFormat() string {
	if len(in.Params) == 0 {
		return ""
	}
	return fmt.Sprintf("%s IN (%s)", in.Field, strings.Join(in.holders, ","))
}

func (in *FieldIN) SQLFormatNotIn() string {
	if len(in.Params) == 0 {
		return ""
	}
	return fmt.Sprintf("%s NOT IN (%s)", in.Field, strings.Join(in.holders, ","))
}

func (in *FieldIN) SQLParams() []interface{} {
	return in.Params
}

//TableJoin
type TableJoin struct {
	tables     map[string]string
	fields     map[string][]string
	from       string
	joins      []*Join
	conditions []string
	orderBy    string
	limit      int
	offset     int
	params     []interface{}
}

type JoinField struct {
	Table string
	Field string
}

type Join struct {
	Method string
	Field1 *JoinField
	Field2 *JoinField
}

func On(tb1 string, f1 string, tb2 string, f2 string) *Join {
	return &Join{
		Method: "INNER",
		Field1: &JoinField{Table: tb1, Field: f1},
		Field2: &JoinField{Table: tb2, Field: f2},
	}
}

func NewTableJoin(table string) *TableJoin {
	tj := &TableJoin{
		tables:     make(map[string]string),
		fields:     make(map[string][]string),
		from:       table,
		joins:      []*Join{},
		conditions: []string{},
		params:     []interface{}{},
	}
	return tj
}

func (tj *TableJoin) Field(table string, fields ...string) {
	if fds, ok := tj.fields[table]; ok {
		fds = append(fds, fields...)
	} else {
		tj.fields[table] = fields
	}
}

func (tj *TableJoin) TableAlias(table string, alias string) {
	tj.tables[table] = alias
}

func (tj *TableJoin) Inner(j *Join) {
	j.Method = "INNER"
	tj.joins = append(tj.joins, j)
}
func (tj *TableJoin) Left(j *Join) {
	j.Method = "LEFT"
	tj.joins = append(tj.joins, j)
}
func (tj *TableJoin) Right(j *Join) {
	j.Method = "RIGHT"
	tj.joins = append(tj.joins, j)
}

func (tj *TableJoin) Condition(condition string, args ...interface{}) {
	tj.conditions = append(tj.conditions, condition)
	tj.params = append(tj.params, args...)
}

func (tj *TableJoin) OrderBy(orderby string) {
	tj.orderBy = orderby
}

func (tj *TableJoin) Limit(offset int, limit int) {
	if limit != 0 {
		tj.limit = limit
	}
	tj.offset = offset
}

func (tj *TableJoin) as(table string) string {
	if alias, ok := tj.tables[table]; ok {
		return fmt.Sprintf("%s AS %s", table, alias)
	}
	return table
}

func (tj *TableJoin) field(table string, field string) string {
	if alias, ok := tj.tables[table]; ok {
		return fmt.Sprintf("%s.%s", alias, field)
	}
	return fmt.Sprintf("%s.%s", table, field)
}

func (tj *TableJoin) on(j *Join) string {
	return fmt.Sprintf("%s JOIN %s ON %s = %s",
		j.Method,
		tj.as(j.Field2.Table),
		tj.field(j.Field1.Table, j.Field1.Field),
		tj.field(j.Field2.Table, j.Field2.Field),
	)
}

func (tj *TableJoin) SQLFormat() string {
	columns := []string{}
	for tb, fs := range tj.fields {
		for _, fd := range fs {
			columns = append(columns, tj.field(tb, fd))
		}
	}
	from := tj.as(tj.from)
	joins := []string{}
	for _, join := range tj.joins {
		joins = append(joins, tj.on(join))
	}
	qs := []string{
		"SELECT",
		"%s",
		"FROM",
		"%s",
		"%s",
	}
	ps := []interface{}{
		strings.Join(columns, ","),
		from,
		strings.Join(joins, " "),
	}
	if len(tj.conditions) != 0 {
		qs = append(qs, "WHERE")
		qs = append(qs, "%s")
		ps = append(ps, strings.Join(tj.conditions, " AND "))
	}
	qs = append(qs, tj.orderBy)
	if tj.limit > 0 {
		if tj.offset > 0 {
			qs = append(qs, "LIMIT %d, %d")
			ps = append(ps, tj.offset, tj.limit)
		} else {
			qs = append(qs, "LIMIT %d")
			ps = append(ps, tj.limit)
		}
	}
	return fmt.Sprintf(strings.Join(qs, " "), ps...)
}

func (tj *TableJoin) SQLParams() []interface{} {
	return tj.params
}

//PasswordFunc
type PasswordFunc func(string) string

func md5Password(password string) string {
	algo := md5.New()
	if _, err := algo.Write([]byte(password)); err != nil {
		return ""
	}
	return hex.EncodeToString(algo.Sum([]byte("orm")))
}

//Cipher
type Cipher interface {
	Encode(string) string
	Decode(string) string
}

type b64Cipher struct{}

func (c *b64Cipher) Encode(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

func (c *b64Cipher) Decode(src string) string {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(decoded)
}

//DBMgr
type _DBMgr struct {
	db DB
}

func DBMgr(db DB) *_DBMgr {
	return &_DBMgr{db: db}
}

func (m *_DBMgr) TableExist(ctx context.Context, database string, table string) (bool, error) {
	qs := []string{
		"SELECT",
		"count(1)",
		"FROM",
		"INFORMATION_SCHEMA.TABLES",
		"WHERE",
		"TABLE_SCHEMA = ?",
		"AND",
		"TABLE_NAME = ?",
	}
	ps := []interface{}{
		database,
		table,
	}

	rows, err := m.db.QueryContext(ctx, strings.Join(qs, " "), ps...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		break
	}
	return count != 0, nil
}

//utils
var Password PasswordFunc
var Encoder Cipher

func init() {
	Password = md5Password
	Encoder = &b64Cipher{}
}
