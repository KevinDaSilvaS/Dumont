package database

import (
	"database/sql"
	"dumont/config"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
}

// Get a database handle. https://go.dev/doc/tutorial/database-access
func Connect(config config.Config) Database {
	cfg := mysql.NewConfig()
	cfg.User = config.User
	cfg.Passwd = config.Passwd
	cfg.Net = "tcp"
	cfg.Addr = config.Host + ":" + config.Port
	cfg.DBName = config.DatabaseName

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
