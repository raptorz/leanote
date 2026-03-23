package postgres

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/leanote/leanote/app/db"
)

func (p *PostgresDatabase) Insert(table string, data interface{}) bool {
	val := reflect.ValueOf(data)
	typVal := val.Type()

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typVal = typVal.Elem()
	}

	columns := make([]string, 0, typVal.NumField())
	values := make([]interface{}, 0, typVal.NumField())
	placeholders := make([]string, 0, typVal.NumField())

	for i := 0; i < typVal.NumField(); i++ {
		field := typVal.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = field.Name
		}
		if tag == "-" {
			continue
		}

		fieldVal := val.Field(i)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}

		columns = append(columns, strings.ToLower(tag))
		values = append(values, fieldVal.Interface())
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(values)))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	var id string
	err := p.db.QueryRow(query, values...).Scan(&id)
	if err != nil {
		return false
	}
	return true
}

func (p *PostgresDatabase) Update(table string, id string, data interface{}) bool {
	val := reflect.ValueOf(data)
	typVal := val.Type()

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typVal = typVal.Elem()
	}

	setParts := make([]string, 0, typVal.NumField())
	values := make([]interface{}, 0, typVal.NumField()+1)

	for i := 0; i < typVal.NumField(); i++ {
		field := typVal.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			tag = field.Name
		}
		if tag == "-" {
			continue
		}

		setParts = append(setParts, fmt.Sprintf("%s = $%d", strings.ToLower(tag), len(values)+1))
		values = append(values, val.Field(i).Interface())
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		table, strings.Join(setParts, ", "), len(values)+1)
	values = append(values, id)

	_, err := p.db.Exec(query, values...)
	if err != nil {
		return false
	}
	return true
}

func (p *PostgresDatabase) Delete(table string, id string) bool {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := p.db.Exec(query, id)
	if err != nil {
		return false
	}
	return true
}

func (p *PostgresDatabase) InsertOld(collection interface{}, data interface{}) bool {
	return true
}

func (p *PostgresDatabase) UpdateByIdAndUserId(collection interface{}, id, userId string, data interface{}) bool {
	return true
}

func (p *PostgresDatabase) GetByIdAndUserId(collection interface{}, id, userId string, data interface{}) {
}

func (p *PostgresDatabase) GetByQ(collection interface{}, query interface{}, data interface{}) {
}

func (p *PostgresDatabase) ListByQ(collection interface{}, query interface{}, data interface{}) {
}

func (p *PostgresDatabase) Count(collection interface{}, query interface{}) int {
	return 0
}

func (p *PostgresDatabase) Has(collection interface{}, query interface{}) bool {
	return false
}

func (p *PostgresDatabase) BatchInsert(table string, data []interface{}) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}

	successCount := 0
	for _, item := range data {
		if p.Insert(table, item) {
			successCount++
		}
	}
	return successCount, nil
}

func (p *PostgresDatabase) BatchUpdate(table string, ids []string, data interface{}) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	successCount := 0
	for _, id := range ids {
		if p.Update(table, id, data) {
			successCount++
		}
	}
	return successCount, nil
}

func (p *PostgresDatabase) BatchDelete(table string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	successCount := 0
	for _, id := range ids {
		if p.Delete(table, id) {
			successCount++
		}
	}
	return successCount, nil
}

func (p *PostgresDatabase) Select(table string, fields []string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	selectClause := "*"
	if len(fields) > 0 {
		selectClause = strings.Join(fields, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s", selectClause, table)
	if where != "" {
		query += fmt.Sprintf(" WHERE %s", where)
	}

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (p *PostgresDatabase) Join(mainTable, joinTable, joinType, onCondition string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s %s %s ON %s", mainTable, joinType, joinTable, onCondition)
	if where != "" {
		query += fmt.Sprintf(" WHERE %s", where)
	}

	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (p *PostgresDatabase) Paginate(table string, page, pageSize int, where string, orderBy string, args ...interface{}) (db.PaginationResult, error) {
	result := db.PaginationResult{
		Page:     page,
		PageSize: pageSize,
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if where != "" {
		countQuery += fmt.Sprintf(" WHERE %s", where)
	}

	var total int64
	err := p.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return result, err
	}

	result.Total = total
	result.TotalPages = int((total + int64(pageSize) - 1) / int64(pageSize))

	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf("SELECT * FROM %s", table)
	if where != "" {
		dataQuery += fmt.Sprintf(" WHERE %s", where)
	}
	if orderBy != "" {
		dataQuery += fmt.Sprintf(" ORDER BY %s", orderBy)
	}
	dataQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := p.db.Query(dataQuery, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		result.Data = append(result.Data, row)
	}

	return result, nil
}

func (p *PostgresDatabase) Begin() (db.Transaction, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	return NewPostgresTransaction(tx), nil
}
