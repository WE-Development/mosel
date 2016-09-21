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
	q["tableExists"] = "SHOW TABLES LIKE "
	q["createNodes"] =
		`CREATE TABLE Nodes (
  			ID   INT          NOT NULL AUTO_INCREMENT,
  			Name VARCHAR(255) NOT NULL,
  			Url  VARCHAR(2000),

  			PRIMARY KEY (ID)
		);`
	q["createDiagrams"] =
		`CREATE TABLE Diagrams (
			  ID   INT          NOT NULL AUTO_INCREMENT,
			  Name VARCHAR(255) NOT NULL,
			  Node INT          NOT NULL,

			  PRIMARY KEY (ID),
			  FOREIGN KEY (Node) REFERENCES Nodes (ID)
		);`

	q["createGraphs"] =
		`CREATE TABLE Graphs (
			  ID      INT          NOT NULL AUTO_INCREMENT,
			  Name    VARCHAR(255) NOT NULL,
			  Diagram INT          NOT NULL,

			  PRIMARY KEY (ID),
			  FOREIGN KEY (Diagram) REFERENCES Diagrams (ID)
		);`

	q["createDataPoints"] =
		`CREATE TABLE DataPoints (
  			ID        INT NOT NULL AUTO_INCREMENT,
  			Value     DOUBLE,
  			Timestamp TIMESTAMP    DEFAULT NOW(),
  			Graph     INT NOT NULL,

			PRIMARY KEY (ID),
  			FOREIGN KEY (Graph) REFERENCES Graphs (ID)
		);`

	q["all"] = "SELECT * FROM bla"
	return q
}