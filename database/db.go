package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
}

type Path struct {
	Variable_name string
	Value         string
}

// Get a database handle. https://go.dev/doc/tutorial/database-access
func Connect() Database {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "oi"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "database_name"

	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return Database{Connection: db}
}

func (db *Database) GetPath() (string, error) {
	var variable Path
	row := db.Connection.QueryRow("SHOW VARIABLES WHERE Variable_name = 'log_bin_basename';")
	if err := row.Scan(&variable.Variable_name, &variable.Value); err != nil {
		return "", fmt.Errorf("getBinLogPath %v", err)
	}
	return variable.ParseBinLogPath(), nil
}

func (p Path) ParseBinLogPath() string {
	parts := strings.Split(p.Value, "/")
	var basePath strings.Builder
	for i := 0; i < len(parts)-1; i++ {
		basePath.WriteString(parts[i] + "/")
	}

	return basePath.String()
}

type BinLogFile struct {
	Log_name  string
	File_size int
}

func (db *Database) GetBinLogFiles() ([]string, error) {
	files := make([]string, 0)

	rows, err := db.Connection.Query("SHOW BINARY LOGS;")
	if err != nil {
		return nil, fmt.Errorf("getBinLogFiles %v", err)
	}

	for rows.Next() {
		var file BinLogFile
		if err := rows.Scan(&file.Log_name, &file.File_size); err != nil {
			continue
		}

		files = append(files, file.Log_name)
	}
	return files, nil
}
