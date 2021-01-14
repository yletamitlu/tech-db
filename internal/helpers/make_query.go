package helpers

import (
	"fmt"
)

func MakeQuery(values []interface{}, query string, limit int, desc bool, since string) (string, []interface{}) {
	i := len(values) + 1

	if since != "" {
		if desc {
			query += fmt.Sprintf(" and id < $%d", i)
		} else {
			query += fmt.Sprintf(" and id > $%d", i)
		}
		i++
		values = append(values, since)
	}

	query += " order by created_at"

	if desc {
		query += " desc"
	}

	query += ", id"

	if desc {
		query += " desc"
	}

	query += " limit " + fmt.Sprintf("$%d", i)
	values = append(values, limit)

	return query, values
}

func MakeQueryForUsers(values []interface{}, query string, limit int, desc bool, since string) (string, []interface{}) {
	i := len(values) + 1

	if since != "" {
		if desc {
			query += fmt.Sprintf(" and nickname < $%d", i)
		} else {
			query += fmt.Sprintf(" and nickname > $%d", i)
		}
		i++
		values = append(values, since)
	}

	query += " order by nickname"

	if desc {
		query += " desc"
	}

	if limit > 0 {
		query += " limit " + fmt.Sprintf("$%d", i)
		values = append(values, limit)
	}

	return query, values
}
