package common

import (
	"fmt"
	"strings"
)

type QueryBuilder struct {
	selectFields []string
	fromTable    string
	joins        []joinClause
	whereClause  string
	whereArgs    []interface{}
	orderBy      string
	limit        int
	offset       int
}

type joinClause struct {
	table       string
	joinType    string
	onCondition string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selectFields = fields
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.fromTable = table
	return qb
}

func (qb *QueryBuilder) Join(table, onCondition string) *QueryBuilder {
	qb.joins = append(qb.joins, joinClause{table: table, joinType: "INNER JOIN", onCondition: onCondition})
	return qb
}

func (qb *QueryBuilder) LeftJoin(table, onCondition string) *QueryBuilder {
	qb.joins = append(qb.joins, joinClause{table: table, joinType: "LEFT JOIN", onCondition: onCondition})
	return qb
}

func (qb *QueryBuilder) Where(clause string, args ...interface{}) *QueryBuilder {
	qb.whereClause = clause
	qb.whereArgs = args
	return qb
}

func (qb *QueryBuilder) OrderBy(order string) *QueryBuilder {
	qb.orderBy = order
	return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.limit = n
	return qb
}

func (qb *QueryBuilder) Offset(n int) *QueryBuilder {
	qb.offset = n
	return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	var sql string
	var args []interface{}

	selectClause := "*"
	if len(qb.selectFields) > 0 {
		selectClause = strings.Join(qb.selectFields, ", ")
	}

	sql = fmt.Sprintf("SELECT %s FROM %s", selectClause, qb.fromTable)

	for _, join := range qb.joins {
		sql += fmt.Sprintf(" %s %s ON %s", join.joinType, join.table, join.onCondition)
	}

	if qb.whereClause != "" {
		sql += fmt.Sprintf(" WHERE %s", qb.whereClause)
		args = append(args, qb.whereArgs...)
	}

	if qb.orderBy != "" {
		sql += fmt.Sprintf(" ORDER BY %s", qb.orderBy)
	}

	if qb.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", qb.limit)
	}

	if qb.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	return sql, args
}
