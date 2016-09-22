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
	q["nodeByName"] = "SELECT * FROM Nodes WHERE Name=?"
	q["diagramByName"] = "SELECT * FROM Diagrams WHERE Name=?"
	q["graphByName"] = "SELECT * FROM Graphs WHERE Name=?"

	q["insertNode"] = "INSERT INTO Nodes (Name, Url) VALUES(?,?)"
	q["insertDiagram"] = "INSERT INTO Diagrams (Name, Node) VALUES (?,?)"
	q["insertGraph"] = "INSERT INTO Graphs (Name, Diagram) VALUES (?,?)"
	q["insertDataPoint"] = "INSERT INTO DataPoints (Value, Timestamp, Graph) VALUES (?,?,?)"

	q["createNodes"] =
		`CREATE TABLE Nodes (
  			ID   INT          NOT NULL AUTO_INCREMENT,
  			Name VARCHAR(255) NOT NULL,
  			Url  VARCHAR(2000),

  			PRIMARY KEY (ID)
		)`
	q["createDiagrams"] =
		`CREATE TABLE Diagrams (
			  ID   INT          NOT NULL AUTO_INCREMENT,
			  Name VARCHAR(255) NOT NULL,
			  Node INT          NOT NULL,

			  PRIMARY KEY (ID),
			  FOREIGN KEY (Node) REFERENCES Nodes (ID)
		)`

	q["createGraphs"] =
		`CREATE TABLE Graphs (
			  ID      INT          NOT NULL AUTO_INCREMENT,
			  Name    VARCHAR(255) NOT NULL,
			  Diagram INT          NOT NULL,

			  PRIMARY KEY (ID),
			  FOREIGN KEY (Diagram) REFERENCES Diagrams (ID)
		)`

	q["createDataPoints"] =
		`CREATE TABLE DataPoints (
  			ID        INT NOT NULL AUTO_INCREMENT,
  			Value     DOUBLE,
  			Timestamp TIMESTAMP    DEFAULT NOW(),
  			Graph     INT NOT NULL,

			PRIMARY KEY (ID),
  			FOREIGN KEY (Graph) REFERENCES Graphs (ID)
		)`

	q["all"] =
		`SELECT
			p.Value "value",
			p.Timestamp "timestamp",
			g.Name "graph",
			d.Name "diagram",
			n.Name "node",
			n.Url "url"
		FROM DataPoints p
			RIGHT JOIN Graphs g
				ON p.Graph=g.Name
			RIGHT JOIN Diagrams d
				ON g.Diagram=d.Name
			RIGHT JOIN Nodes n
				ON d.Node=n.Name`

	return q
}