package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/x-mod/errors"
)

var (
	_ sql.DB
	_ time.Time
	_ strings.Reader
)

type User struct {
	Id       int32      `json:"id" yaml:"id"`
	Name     string     `json:"name" yaml:"name"`
	Age      int32      `json:"age" yaml:"age"`
	Sex      bool       `json:"sex" yaml:"sex"`
	FooBar   int32      `json:"foo" yaml:"bar"`
	CreateAt time.Time  `json:"createAt" yaml:"createAt"`
	UpdateAt *time.Time `json:"updateAt" yaml:"updateAt"`
}

var UserDBColumns = []string{
	"`id`",
	"`name`",
	"`age`",
	"`sex`",
	"`foo_bar`",
	"`create_at`",
	"`update_at`",
}

var UserDBTable = "user"

type _UserMgr struct {
}

var UserMgr *_UserMgr

func (m *_UserMgr) NewUser() *User {
	return &User{}
}

type _UserDBMgr struct {
	db DB
}

func UserDBMgr(db DB) *_UserDBMgr {
	return &_UserDBMgr{db: db}
}

type UserPK struct {
	Id int32
}

func (m *_UserMgr) NewPK() *UserPK {
	return &UserPK{}
}

func (u *UserPK) SQLFormat() []string {
	return []string{
		"`id` = ?",
	}
}

func (u *UserPK) SQLParams() []interface{} {
	return []interface{}{
		u.Id,
	}
}

func (u *UserPK) DBColumns() []string {
	return []string{
		"`id`",
	}
}

type UserNameUK struct {
	Name string
}

func (m *_UserMgr) NewUserNameUK() *UserNameUK {
	return &UserNameUK{}
}

func (u *UserNameUK) SQLFormat() []string {
	return []string{
		"`name` = ?",
	}
}

func (u *UserNameUK) SQLParams() []interface{} {
	return []interface{}{
		u.Name,
	}
}

func (u *UserNameUK) DBColumns() []string {
	return []string{
		"`name`",
	}
}

type UserAgeSexIndex struct {
	Age int32
	Sex bool
}

func (m *_UserMgr) NewUserAgeSexIndex() *UserAgeSexIndex {
	return &UserAgeSexIndex{}
}

func (u *UserAgeSexIndex) SQLFormat() []string {
	return []string{
		"`age` = ?",
		"`sex` = ?",
	}
}

func (u *UserAgeSexIndex) SQLParams() []interface{} {
	return []interface{}{
		u.Age,
		u.Sex,
	}
}

func (u *UserAgeSexIndex) DBColumns() []string {
	return []string{
		"`age`",
		"`sex`",
	}
}

func (m *_UserDBMgr) FetchByPK(ctx context.Context, pk *UserPK) (*User, error) {
	return m.FetchOne(ctx, pk.SQLFormat(), pk.SQLParams()...)
}

func (m *_UserDBMgr) FetchOne(ctx context.Context, conditions []string, args ...interface{}) (*User, error) {
	objs, err := m.QueryContext(ctx, conditions, "", 0, 1, args...)
	if err != nil {
		return nil, err
	}
	if len(objs) == 0 {
		return nil, errors.New("no record")
	}
	return objs[0], nil
}

func (m *_UserDBMgr) Fetch(ctx context.Context, conditions []string, args ...interface{}) ([]*User, error) {
	return m.QueryContext(ctx, conditions, "", 0, 0, args...)
}

func (m *_UserDBMgr) CountContext(ctx context.Context, conditions []string, args ...interface{}) (int64, error) {
	qs := []string{
		"SELECT",
		"count(1)",
		"FROM",
		"%s",
	}
	ps := []interface{}{
		"user",
	}
	if len(conditions) != 0 {
		qs = append(qs, "WHERE")
		qs = append(qs, "%s")
		ps = append(ps, strings.Join(conditions, " AND "))
	}
	q := fmt.Sprintf(strings.Join(qs, " "), ps...)
	return m.countContext(ctx, q, args...)
}

func (m *_UserDBMgr) QueryContext(ctx context.Context, conditions []string, orderBy string, offset int, limit int, args ...interface{}) ([]*User, error) {
	qs := []string{
		"SELECT",
		"%s",
		"FROM",
		"%s",
	}
	ps := []interface{}{
		strings.Join(UserDBColumns, ","),
		"user",
	}
	if len(conditions) != 0 {
		qs = append(qs, "WHERE")
		qs = append(qs, "%s")
		ps = append(ps, strings.Join(conditions, " AND "))
	}
	qs = append(qs, orderBy)
	if limit > 0 {
		if offset > 0 {
			qs = append(qs, "LIMIT %d, %d")
			ps = append(ps, offset, limit)
		} else {
			qs = append(qs, "LIMIT %d")
			ps = append(ps, limit)
		}
	}
	q := fmt.Sprintf(strings.Join(qs, " "), ps...)
	return m.queryContext(ctx, q, args...)
}

func (m *_UserDBMgr) queryContext(ctx context.Context, q string, args ...interface{}) ([]*User, error) {
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, errors.Annotatef(err, "query %s %v", q, args)
	}
	defer rows.Close()

	objs := []*User{}
	for rows.Next() {
		var obj User
		var Age sql.NullInt32
		var FooBar sql.NullInt32
		var CreateAt int64
		var UpdateAt int64

		err = rows.Scan(&(obj.Id), &(obj.Name),
			&Age, &(obj.Sex),
			&FooBar,
			&CreateAt,
			&UpdateAt)
		if err != nil {
			return nil, err
		}

		if Age.Valid {
			obj.Age = int32(Age.Int32)
		}
		if FooBar.Valid {
			obj.FooBar = int32(FooBar.Int32)
		}
		tmCreateAt := time.Unix(CreateAt, 0)
		obj.CreateAt = tmCreateAt
		tmUpdateAt := time.Unix(UpdateAt, 0)
		if UpdateAt > 0 {
			obj.UpdateAt = &tmUpdateAt
		}

		objs = append(objs, &obj)
	}
	return objs, nil
}

func (m *_UserDBMgr) countContext(ctx context.Context, q string, args ...interface{}) (int64, error) {
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return 0, errors.Annotatef(err, "query %s %v", q, args)
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
		break
	}
	return count, nil
}

func (m *_UserDBMgr) DeleteByPK(ctx context.Context, pk *UserPK) (int64, error) {
	if pk == nil {
		return 0, errors.New("delete pk required")
	}
	return m.Delete(ctx, pk.SQLFormat(), pk.SQLParams()...)
}

func (m *_UserDBMgr) Delete(ctx context.Context, conditions []string, args ...interface{}) (int64, error) {
	if len(conditions) > 0 {
		qs := []string{
			"DELETE",
			"FROM",
			"%s",
			"WHERE",
			"%s",
		}
		ps := []interface{}{
			"user",
			strings.Join(conditions, " AND "),
		}
		q := fmt.Sprintf(strings.Join(qs, " "), ps...)
		return m.execContext(ctx, q, args...)
	}
	return 0, errors.New("delete conditions required")
}

func (m *_UserDBMgr) execContext(ctx context.Context, q string, args ...interface{}) (int64, error) {
	result, err := m.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, errors.Annotatef(err, "exec %s %v", q, args)
	}
	return result.RowsAffected()
}

func (m *_UserDBMgr) UpdateByPK(ctx context.Context, updates map[string]interface{}, pk *UserPK) (int64, error) {
	return m.Update(ctx, updates, pk.SQLFormat(), pk.SQLParams()...)
}

func (m *_UserDBMgr) Update(ctx context.Context, updates map[string]interface{}, conditions []string, args ...interface{}) (int64, error) {
	qs := []string{
		"UPDATE",
		"%s",
		"SET",
		"%s",
		"WHERE",
		"%s",
	}
	ps := []interface{}{
		"user",
	}
	fields := []string{}
	values := []interface{}{}
	for k, v := range updates {
		fields = append(fields, fmt.Sprintf("%s = ?", k))
		values = append(values, v)
	}
	ps = append(ps, strings.Join(fields, ","))
	ps = append(ps, strings.Join(conditions, " AND "))
	values = append(values, args...)
	q := fmt.Sprintf(strings.Join(qs, " "), ps...)
	return m.execContext(ctx, q, values...)
}

func (m *_UserDBMgr) Create(ctx context.Context, obj *User) (int64, error) {
	columns := []string{
		"`name`",
		"`age`",
		"`sex`",
		"`foo_bar`",
		"`create_at`",
		"`update_at`",
	}
	qs := []string{
		"INSERT",
		"INTO",
		"%s",
		"(%s)",
		"VALUES",
		"%s",
	}
	ps := []interface{}{
		"user",
		strings.Join(columns, ","),
	}
	blanks := []string{}
	values := []interface{}{}

	blank := []string{}
	for i := 0; i < len(columns); i++ {
		blank = append(blank, "?")
	}
	blanks = append(blanks, fmt.Sprintf("(%s)", strings.Join(blank, ",")))
	values = append(values, obj.Name)
	values = append(values, obj.Age)
	values = append(values, obj.Sex)
	values = append(values, obj.FooBar)
	values = append(values, obj.CreateAt.Unix())
	if obj.UpdateAt != nil {
		values = append(values, obj.UpdateAt.Unix())
	} else {
		values = append(values, 0)
	}

	ps = append(ps, strings.Join(blanks, ","))
	q := fmt.Sprintf(strings.Join(qs, " "), ps...)

	result, err := m.db.ExecContext(ctx, q, values...)
	if err != nil {
		return 0, errors.Annotatef(err, "exec %s %v", q, values)
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	obj.Id = int32(lastInsertId)
	return result.RowsAffected()
}

func (m *_UserDBMgr) BatchCreate(ctx context.Context, objs []*User) (int64, error) {
	columns := []string{
		"`name`",
		"`age`",
		"`sex`",
		"`foo_bar`",
		"`create_at`",
		"`update_at`",
	}
	qs := []string{
		"INSERT",
		"INTO",
		"%s",
		"(%s)",
		"VALUES",
		"%s",
	}
	ps := []interface{}{
		"user",
		strings.Join(columns, ","),
	}
	blanks := []string{}
	values := []interface{}{}
	for _, obj := range objs {
		blank := []string{}
		for i := 0; i < len(columns); i++ {
			blank = append(blank, "?")
		}
		blanks = append(blanks, fmt.Sprintf("(%s)", strings.Join(blank, ",")))
		values = append(values, obj.Name)
		values = append(values, obj.Age)
		values = append(values, obj.Sex)
		values = append(values, obj.FooBar)
		values = append(values, obj.CreateAt.Unix())
		if obj.UpdateAt != nil {
			values = append(values, obj.UpdateAt.Unix())
		} else {
			values = append(values, 0)
		}
	}
	ps = append(ps, strings.Join(blanks, ","))
	q := fmt.Sprintf(strings.Join(qs, " "), ps...)
	return m.execContext(ctx, q, values...)
}
