package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/x-mod/errors"
)

var (
	_ sql.DB
	_ time.Time
)

type Summary struct {
	Id   int32  `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Age  int32  `json:"age" yaml:"age"`
	Sex  bool   `json:"sex" yaml:"sex"`
}

var SummaryDBColumns = []string{
	"`id`",
	"`name`",
	"`age`",
	"`sex`",
}

type _SummaryMgr struct {
}

var SummaryMgr *_SummaryMgr

func (m *_SummaryMgr) NewSummary() *Summary {
	return &Summary{}
}

type _SummaryDBMgr struct {
	db DB
}

func SummaryDBMgr(db DB) *_SummaryDBMgr {
	return &_SummaryDBMgr{db: db}
}

func (m *_SummaryDBMgr) QueryContext(ctx context.Context, q string, args ...interface{}) ([]*Summary, error) {
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, errors.Annotatef(err, "query %s %v", q, args)
	}
	defer rows.Close()

	objs := []*Summary{}
	for rows.Next() {
		var obj Summary
		var Age sql.NullInt32

		err = rows.Scan(&(obj.Id), &(obj.Name),
			&Age, &(obj.Sex))
		if err != nil {
			return nil, err
		}

		if Age.Valid {
			obj.Age = int32(Age.Int32)
		}

		objs = append(objs, &obj)
	}
	return objs, nil
}

func (m *_SummaryDBMgr) CountContext(ctx context.Context, q string, args ...interface{}) (int64, error) {
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
