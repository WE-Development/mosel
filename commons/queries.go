package commons

import "fmt"

type SqlQueries map[string]string

func GetQueries(dialect string) (SqlQueries, error) {
	if (dialect == "mysql") {
		return getMysqlQueries(), nil
	}

	return nil, fmt.Errorf("SQL dialect %s is not supported", dialect)
}

func getMysqlQueries() SqlQueries {
	q := make(SqlQueries)
	q["all"] = "SELECT * FROM bla"
	return q
}