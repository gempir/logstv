package main

import (
	"fmt"
	"strings"
)

type order string

var (
	orderDesc order = "DESC"
	orderAsc  order = "ASC"
)

func buildQuery(selectFields []string, table string, whereClauses []string, orderBy order, limit int) string {
	selectFieldsJoined := strings.Join(selectFields, ",")
	whereClausesJoined := strings.Join(whereClauses, " AND ")

	var limitString string
	if limit > 0 {
		limitString = fmt.Sprintf("LIMIT %d", limit)
	}

	return fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s ORDER BY timestamp %s %s",
		selectFieldsJoined,
		table,
		whereClausesJoined,
		orderBy,
		limitString,
	)
}
